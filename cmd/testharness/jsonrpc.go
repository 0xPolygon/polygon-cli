package testharness

// junkJSONRPC is a static list of various test strings some of which are valid json but invalid json rpc. Some are just complete trash
var junkJSONRPC = []string{
	// These should be acceptable
	`{"jsonrpc": "2.0", "result": null, "id": 1}`,                                             // null result should be fine
	`{"jsonrpc": "2.0", "result": 123, "id": 1}`,                                              // number should be fine
	`{"jsonrpc": "2.0", "result": "123", "id": 1}`,                                            // string should be fine
	`{"jsonrpc": "2.0", "result": [1, 2 ,3], "id": 1}`,                                        // array should be fine
	`{"jsonrpc": "2.0", "result": {}, "id": 1}`,                                               // empty object should be okay
	`{"jsonrpc": "2.0", "result": {"1": true, "2": true, "3": true}, "id": 1}`,                // object should be okay
	`{"jsonrpc": "2.0", "result": 123, "id": "abc123"}`,                                       // string id should be okay
	`[{"jsonrpc": "2.0", "result": 123, "id": 1},{"jsonrpc": "2.0", "result": 456, "id": 2}]`, // batch response
	`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "test err"}, "id": null}`,        // error response should be fine

	// these are dubious
	`{"jsonrpc": "2.0", "id": 1}`,                                                                        // completely missing result
	`{"jsonrpc": "2", "result": 123, "id": 1}`,                                                           // 2 instead of 2.0
	`{"jsonrpc": "2.0", "result": 123, "id": null}`,                                                      // null id
	`{"jsonrpc": "2.0", "result": 123, "id": 1}` + string([]byte{0}),                                     // trailing null byte
	`{"jsonrpc": "2.0", "result": 123, "id": 1}` + "\r\n" + `{"jsonrpc": "2.0", "result": 123, "id": 1}`, // \r\n seperator
	`{"jsonrpc": "2.0", "result": 123, "id": 1}` + "\r" + `{"jsonrpc": "2.0", "result": 123, "id": 1}`,   // \r seperator
	`{"jsonrpc": "2.0", "result": 123, "id": 1}` + "\n" + `{"jsonrpc": "2.0", "result": 123, "id": 1}`,   // \n seperator
	`{"jsonrpc": "2.0", "result": 123}`,                                                                  // missing id all together
	`{"result": 123, "id": 1}`,                                                                           // missing json rpc all together.. This would be jsonrpc 1.0
	`[{"jsonrpc": "2.0", "result": 123, "id": 1},{"jsonrpc": "2.0", "result": 123, "id": 1}]`,            // batch response with two responses fro the same id
	`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "test err"}, "id": 1}`,                      // error id must be null
	`{"jsonrpc": "2.0", "error": {"code": "foo", "message": "test err"}, "id": null}`,
	`{"jsonrpc": "2.0", "error": {"code": 2.718282828, "message": "test err"}, "id": null}`, // non-integer code
	`{"jsonrpc": "2.0", "error": {}, "id": null}`,                                           // empty object for error
	`{"jsonrpc": "2.0", "error": null, "id": null}`,                                         // numm error
	`{"jsonrpc": "2.0", "error": {"message": "test err"}, "id": null}`,                      // Missing code
	`{"jsonrpc": "2.0", "error": {"code": -32700}, "id": null}`,                             // missing message
	`{"jsonrpc": "2.0", "error": true, "id": null}`,                                         // error is wrong type

	// these should break something
	`{"jsonrpc": "2.0", "result": , "id": 1}`,                   // missing result ... broken json
	`{"jsonrpc": "2.0", "result": nil, "id": 1}`,                // nil instead of null
	`{"jsonrpc": "2.0", "result": 123, "result": 456, "id": 1}`, // result specified twice? might work
	``,                      // valid json but probably will break
	`null`,                  // valid but should break
	`0`,                     // valid but should break rp
	`0x00`,                  // invalid I think
	`<xml />`,               // wrong type
	`hi`,                    // invalid json
	`"hi"`,                  // valid json but should break
	string([]byte{0, 0, 0}), // null bytes should break
}

var junkContentTypeHeader = []string{
	"application/javascript",
	"application/octet-stream",
	"application/xhtml+xml",
	"application/json",
	"application/xml",
	"application/x-www-form-urlencoded",
	"audio/x-wav",
	"image/gif",
	"image/png",
	"multipart/mixed",
	"multipart/form-data",
}
