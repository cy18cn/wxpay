package wxpay

import (
	"bytes"
	"encoding/xml"
	"io"
)

// UnmarshalToMap unmarshal xml into map
func UnmarshalToMap(b []byte) map[string]interface{} {
	decoder := xml.NewDecoder(bytes.NewBuffer(b))
	s := NewStack()
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		switch ele := t.(type) {
		case xml.StartElement:
			s.Push(ele)
		case xml.CharData:
			s.Push(string(xml.CharData(ele)))
		case xml.EndElement:
			entry := make(map[string]interface{})
			var val interface{}
			for {
				p := s.Pop()

				// push key [start] val into stack
				if start, ok := p.(xml.StartElement); ok {
					node := make(map[string]interface{})
					node[start.Name.Local] = val
					s.Push(node)
					break
				}

				switch p.(type) {
				case map[string]interface{}:
					// handle pushed map
					node := p.(map[string]interface{})
					for k, v := range node {
						entry[k] = v
					}
					if val == nil {
						val = entry
					}
				default:
					// handle xml char data
					val = p
				}
				// handle pushed map
				//if node, ok := p.(map[string]interface{}); ok {
				//	for k, v := range node {
				//		entry[k] = v
				//	}
				//	if val == nil {
				//		val = entry
				//	}
				//	continue
				//}
				//
				//// handle xml char data
				//val = p
			}
		}

	}

	return s.Pop().(map[string]interface{})
}
