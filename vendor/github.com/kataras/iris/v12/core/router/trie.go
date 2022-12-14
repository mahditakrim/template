package router

import (
	"strings"

	"github.com/kataras/iris/v12/context"
)

const (
	// ParamStart the character in string representation where the underline router starts its dynamic named parameter.
	ParamStart = ":"
	// paramStartCharacter is the character as rune of ParamStart.
	paramStartCharacter = ':'
	// WildcardParamStart the character in string representation where the underline router starts its dynamic wildcard
	// path parameter.
	WildcardParamStart = "*"
	// wildcardParamStartCharacter is the character as rune of WildcardParamStart.
	wildcardParamStartCharacter = '*'
)

// An iris-specific identical version of the https://github.com/kataras/muxie version 1.0.0 released at 15 Oct 2018
type trieNode struct {
	parent *trieNode

	children               map[string]*trieNode
	hasDynamicChild        bool     // does one of the children contains a parameter or wildcard?
	childNamedParameter    bool     // is the child a named parameter (single segmnet)
	childWildcardParameter bool     // or it is a wildcard (can be more than one path segments) ?
	paramKeys              []string // the param keys without : or *.
	end                    bool     // it is a complete node, here we stop and we can say that the node is valid.
	key                    string   // if end == true then key is filled with the original value of the insertion's key.
	// if key != "" && its parent has childWildcardParameter == true,
	// we need it to track the static part for the closest-wildcard's parameter storage.
	staticKey string

	// insert data.
	Route    context.RouteReadOnly
	Handlers context.Handlers
}

func newTrieNode() *trieNode {
	n := new(trieNode)
	return n
}

func (tn *trieNode) hasChild(s string) bool {
	return tn.getChild(s) != nil
}

func (tn *trieNode) getChild(s string) *trieNode {
	if tn.children == nil {
		return nil
	}

	return tn.children[s]
}

func (tn *trieNode) addChild(s string, n *trieNode) {
	if tn.children == nil {
		tn.children = make(map[string]*trieNode)
	}

	if _, exists := tn.children[s]; exists {
		return
	}

	n.parent = tn
	tn.children[s] = n
}

func (tn *trieNode) findClosestParentWildcardNode() *trieNode {
	tn = tn.parent
	for tn != nil {
		if tn.childWildcardParameter {
			return tn.getChild(WildcardParamStart)
		}

		tn = tn.parent
	}

	return nil
}

func (tn *trieNode) String() string {
	return tn.key
}

type trie struct {
	root *trieNode

	// if true then it will handle any path if not other parent wildcard exists,
	// so even 404 (on http services) is up to it, see trie#insert.
	hasRootWildcard bool
	hasRootSlash    bool

	statusCode int // for error codes only, method is ignored.
	method     string

	// subdomain is empty for default-hostname routes,
	// ex: mysubdomain.
	subdomain string
}

const (
	pathSep  = "/"
	pathSepB = '/'
)

func slowPathSplit(path string) []string {
	if path == "/" {
		return []string{"/"}
	}

	return strings.Split(path, pathSep)[1:]
}

func (tr *trie) insert(path string, route context.RouteReadOnly, handlers context.Handlers) {
	input := slowPathSplit(path)
	if len(input) == 0 {
		return
	}

	n := tr.root
	if path == pathSep {
		tr.hasRootSlash = true
	}

	var paramKeys []string

	for _, s := range input {
		if len(s) == 0 {
			continue
		}

		c := s[0]

		if len(s) > 1 { // has more than one character.
			// get the next character, should be the name of the parameter.
			// E.g:
			// If /test/:param (or /test/*param) then it's dynamic.
			// If /test/: (or /test/*) then it's static.
			if isParam, isWildcard := c == paramStartCharacter, c == wildcardParamStartCharacter; isParam || isWildcard {
				n.hasDynamicChild = true
				paramKeys = append(paramKeys, s[1:]) // without : or *.
				if isParam {
					n.childNamedParameter = true
					s = ParamStart
				}

				if isWildcard {
					n.childWildcardParameter = true
					s = WildcardParamStart
					if tr.root == n {
						tr.hasRootWildcard = true
					}
				}
			}
		}

		if !n.hasChild(s) {
			child := newTrieNode()
			n.addChild(s, child)
		}

		n = n.getChild(s)
	}

	n.Route = route
	n.Handlers = handlers

	n.paramKeys = paramKeys
	n.key = path
	n.end = true

	i := strings.Index(path, ParamStart)
	if i == -1 {
		i = strings.Index(path, WildcardParamStart)
	}
	if i == -1 {
		i = len(n.key)
	}

	n.staticKey = path[:i]
	// fmt.Printf("trie.insert: (whole path=%v) Path: %s, Route name: %s, Handlers len: %d\n", n.end, n.key, route.Name(), len(handlers))
}

func (tr *trie) search(q string, params *context.RequestParams) *trieNode {
	end := len(q)

	if end == 0 || (end == 1 && q[0] == pathSepB) {
		// fixes only root wildcard but no / registered at.
		if tr.hasRootSlash {
			return tr.root.getChild(pathSep)
		} else if tr.hasRootWildcard {
			// no need to going through setting parameters, this one has not but it is wildcard.
			return tr.root.getChild(WildcardParamStart)
		}

		return nil
	}

	n := tr.root
	start := 1
	i := 1
	var paramValues []string

	for {
		if i == end || q[i] == pathSepB {
			segment := q[start:i]
			if child := n.getChild(segment); child != nil {
				n = child
			} else if n.childNamedParameter {
				n = n.getChild(ParamStart)
				if ln := len(paramValues); cap(paramValues) > ln {
					paramValues = paramValues[:ln+1]
					paramValues[ln] = segment
				} else {
					paramValues = append(paramValues, segment)
				}
			} else if n.childWildcardParameter {
				n = n.getChild(WildcardParamStart)
				if ln := len(paramValues); cap(paramValues) > ln {
					paramValues = paramValues[:ln+1]
					paramValues[ln] = q[start:]
				} else {
					paramValues = append(paramValues, q[start:])
				}
				break
			} else {
				n = n.findClosestParentWildcardNode()
				if n != nil && len(n.paramKeys) > 0 {
					// means that it has :param/static and *wildcard, we go trhough the :param
					// but the next path segment is not the /static, so go back to *wildcard
					// instead of not found.
					//
					// Fixes:
					// /hello/*p
					// /hello/:p1/static/:p2
					// req: http://localhost:8080/hello/dsadsa/static/dsadsa => found
					// req: http://localhost:8080/hello/dsadsa => but not found!
					// and
					// /second/wild/*p
					// /second/wild/static/otherstatic/
					// req: /second/wild/static/otherstatic/random => but not found!
					params.Set(n.paramKeys[0], q[len(n.staticKey):])
					return n
				}

				return nil
			}

			if i == end {
				break
			}

			i++
			start = i
			continue
		}

		i++
	}

	if n == nil || !n.end {
		if n != nil { // we need it on both places, on last segment (below) or on the first unnknown (above).
			if n = n.findClosestParentWildcardNode(); n != nil && len(n.paramKeys) > 0 {
				params.Set(n.paramKeys[0], q[len(n.staticKey):])
				return n
			}
		}

		if tr.hasRootWildcard {
			// that's the case for root wildcard, tests are passing
			// even without it but stick with it for reference.
			// Note ote that something like:
			// Routes: /other2/*myparam and /other2/static
			// Reqs: /other2/staticed will be handled
			// the /other2/*myparam and not the root wildcard, which is what we want.
			//
			n = tr.root.getChild(WildcardParamStart)
			if len(n.paramKeys) == 0 { // fix crashes on /*/*/*.
				return nil
			}

			params.Set(n.paramKeys[0], q[1:])
			return n
		}

		return nil
	}

	for i, paramValue := range paramValues {
		if len(n.paramKeys) > i {
			params.Set(n.paramKeys[i], paramValue)
		}
	}

	return n
}
