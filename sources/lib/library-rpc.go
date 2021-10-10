

package zrun


import "net"
import "net/rpc"
import "strings"
import "sync"




type LibraryRpcServer struct {
	library LibraryStore
	url string
	listener net.Listener
	rpc *rpc.Server
	mutex sync.RWMutex
	waiter sync.WaitGroup
}

type LibraryRpcServerExports struct {
	server *LibraryRpcServer
}




type LibraryRpcClient struct {
	url string
	rpc *rpc.Client
}




func NewLibraryRpcServer (_library LibraryStore, _url string) (*LibraryRpcServer, *Error) {
	var _network, _address string
	if _network_0, _address_0, _error := urlSplit (_url); _error == nil {
		_network = _network_0
		_address = _address_0
	} else {
		return nil, _error
	}
//	logf ('d', 0x49d3cc32, "listening library RPC server on `%s` / `%s`...", _network, _address)
	var _listener net.Listener
	if _listener_0, _error := net.Listen (_network, _address); _error == nil {
		_listener = _listener_0
	} else {
		return nil, errorw (0x565a3b35, _error)
	}
	_exports := & LibraryRpcServerExports {}
	_rpc := rpc.NewServer ()
	if _error := _rpc.RegisterName ("Library", _exports); _error != nil {
		return nil, errorw (0x6c72d486, _error)
	}
	_server := & LibraryRpcServer {
			library : _library,
			url : _url,
			listener : _listener,
			rpc : _rpc,
		}
	_exports.server = _server
	return _server, nil
}


func (_server *LibraryRpcServer) Url () (string) {
	return _server.url
}




func (_server *LibraryRpcServer) Serve () (*Error) {
	_server.waiter.Add (1)
	go _server.loop ()
	_server.waiter.Wait ()
	return nil
}

func (_server *LibraryRpcServer) ServeStart () (*Error) {
	_server.waiter.Add (1)
	go _server.loop ()
	return nil
}

func (_server *LibraryRpcServer) ServeStop () (*Error) {
	if _error := _server.listener.Close (); _error != nil {
		return errorw (0x8748e4e6, _error)
	}
	_server.waiter.Wait ()
	return nil
}


func (_server *LibraryRpcServer) loop () () {
//	logf ('d', 0x205ad5d1, "begin accepting client connections...")
	for {
		if _connection, _error := _server.listener.Accept (); _error == nil {
			_server.waiter.Add (1)
			go _server.handle (_connection)
		} else if strings.HasSuffix (_error.Error (), ": use of closed network connection") {
			break
		} else {
			logError ('w', errorw (0x2737c361, _error))
			break
		}
	}
	_server.waiter.Done ()
//	logf ('d', 0x9ab6240d, "ended accepting client connections;")
}


func (_server *LibraryRpcServer) handle (_connection net.Conn) () {
//	logf ('d', 0x1699394e, "begin handling client connection...")
	_server.rpc.ServeConn (_connection)
//	logf ('d', 0x9f91f3f9, "ended handling client connection;")
	_server.waiter.Done ()
}




func NewLibraryRpcClient (_url string) (*LibraryRpcClient, *Error) {
	var _network, _address string
	if _network_0, _address_0, _error := urlSplit (_url); _error == nil {
		_network = _network_0
		_address = _address_0
	} else {
		return nil, _error
	}
//	logf ('d', 0x1afdd51f, "connecting library RPC client to `%s` / `%s`...", _network, _address)
	var _rpc *rpc.Client
	if _rpc_0, _error := rpc.Dial (_network, _address); _error == nil {
		_rpc = _rpc_0
	} else {
		return nil, errorw (0xf4366bf1, _error)
	}
	_client := & LibraryRpcClient {
			url : _url,
			rpc : _rpc,
		}
	return _client, nil
}


func (_client *LibraryRpcClient) Close () (*Error) {
	if _error := _client.rpc.Close (); _error == nil {
		return nil
	} else {
		return errorw (0x7d41d83f, _error)
	}
}


func (_client *LibraryRpcClient) Url () (string) {
	return _client.url
}




type LibraryRpc_Fingerprint_Input struct {
}

type LibraryRpc_Fingerprint_Output struct {
	Fingerprint string
	Error *Error
}

func (_client *LibraryRpcClient) Fingerprint () (string, *Error) {
	_input := LibraryRpc_Fingerprint_Input {}
	_output := LibraryRpc_Fingerprint_Output {}
	if _error := _client.rpc.Call ("Library.Fingerprint", &_input, &_output); _error == nil {
		return _output.Fingerprint, _output.Error
	} else {
		return "", errorw (0xdbfbf38a, _error)
	}
}

func (_exports *LibraryRpcServerExports) Fingerprint (_input *LibraryRpc_Fingerprint_Input, _output *LibraryRpc_Fingerprint_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Fingerprint, _output.Error = _exports.server.library.Fingerprint ()
	return nil
}




type LibraryRpc_SelectFingerprints_Input struct {
}

type LibraryRpc_SelectFingerprints_Output struct {
	Fingerprints []string
	Error *Error
}

func (_client *LibraryRpcClient) SelectFingerprints () ([]string, *Error) {
	_input := LibraryRpc_SelectFingerprints_Input {}
	_output := LibraryRpc_SelectFingerprints_Output {}
	if _error := _client.rpc.Call ("Library.SelectFingerprints", &_input, &_output); _error == nil {
		return _output.Fingerprints, _output.Error
	} else {
		return nil, errorw (0xf6ef7108, _error)
	}
}

func (_exports *LibraryRpcServerExports) SelectFingerprints (_input *LibraryRpc_SelectFingerprints_Input, _output *LibraryRpc_SelectFingerprints_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Fingerprints, _output.Error = _exports.server.library.SelectFingerprints ()
	return nil
}




type LibraryRpc_ResolveFullByFingerprint_Input struct {
	Fingerprint string
}

type LibraryRpc_ResolveFullByFingerprint_Output struct {
	Scriptlet *Scriptlet
	Error *Error
}

func (_client *LibraryRpcClient) ResolveFullByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	_input := LibraryRpc_ResolveFullByFingerprint_Input {}
	_output := LibraryRpc_ResolveFullByFingerprint_Output {}
	_input.Fingerprint = _fingerprint
	if _error := _client.rpc.Call ("Library.ResolveFullByFingerprint", &_input, &_output); _error == nil {
		return _output.Scriptlet, _output.Error
	} else {
		return nil, errorw (0xfee982f0, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveFullByFingerprint (_input *LibraryRpc_ResolveFullByFingerprint_Input, _output *LibraryRpc_ResolveFullByFingerprint_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Scriptlet, _output.Error = _exports.server.library.ResolveFullByFingerprint (_input.Fingerprint)
	return nil
}




type LibraryRpc_ResolveMetaByFingerprint_Input struct {
	Fingerprint string
}

type LibraryRpc_ResolveMetaByFingerprint_Output struct {
	Scriptlet *Scriptlet
	Error *Error
}

func (_client *LibraryRpcClient) ResolveMetaByFingerprint (_fingerprint string) (*Scriptlet, *Error) {
	_input := LibraryRpc_ResolveMetaByFingerprint_Input {}
	_output := LibraryRpc_ResolveMetaByFingerprint_Output {}
	_input.Fingerprint = _fingerprint
	if _error := _client.rpc.Call ("Library.ResolveMetaByFingerprint", &_input, &_output); _error == nil {
		return _output.Scriptlet, _output.Error
	} else {
		return nil, errorw (0x58a5846f, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveMetaByFingerprint (_input *LibraryRpc_ResolveMetaByFingerprint_Input, _output *LibraryRpc_ResolveMetaByFingerprint_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Scriptlet, _output.Error = _exports.server.library.ResolveMetaByFingerprint (_input.Fingerprint)
	return nil
}




type LibraryRpc_ResolveBodyByFingerprint_Input struct {
	Fingerprint string
}

type LibraryRpc_ResolveBodyByFingerprint_Output struct {
	Body string
	Found bool
	Error *Error
}

func (_client *LibraryRpcClient) ResolveBodyByFingerprint (_fingerprint string) (string, bool, *Error) {
	_input := LibraryRpc_ResolveBodyByFingerprint_Input {}
	_output := LibraryRpc_ResolveBodyByFingerprint_Output {}
	_input.Fingerprint = _fingerprint
	if _error := _client.rpc.Call ("Library.ResolveBodyByFingerprint", &_input, &_output); _error == nil {
		return _output.Body, _output.Found, _output.Error
	} else {
		return "", false, errorw (0x612519d7, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveBodyByFingerprint (_input *LibraryRpc_ResolveBodyByFingerprint_Input, _output *LibraryRpc_ResolveBodyByFingerprint_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Body, _output.Found, _output.Error = _exports.server.library.ResolveBodyByFingerprint (_input.Fingerprint)
	return nil
}




type LibraryRpc_SelectLabels_Input struct {
}

type LibraryRpc_SelectLabels_Output struct {
	Labels []string
	Error *Error
}

func (_client *LibraryRpcClient) SelectLabels () ([]string, *Error) {
	_input := LibraryRpc_SelectLabels_Input {}
	_output := LibraryRpc_SelectLabels_Output {}
	if _error := _client.rpc.Call ("Library.SelectLabels", &_input, &_output); _error == nil {
		return _output.Labels, _output.Error
	} else {
		return nil, errorw (0x149f3d3e, _error)
	}
}

func (_exports *LibraryRpcServerExports) SelectLabels (_input *LibraryRpc_SelectLabels_Input, _output *LibraryRpc_SelectLabels_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Labels, _output.Error = _exports.server.library.SelectLabels ()
	return nil
}




type LibraryRpc_SelectLabelsAll_Input struct {
}

type LibraryRpc_SelectLabelsAll_Output struct {
	Labels []string
	Error *Error
}

func (_client *LibraryRpcClient) SelectLabelsAll () ([]string, *Error) {
	_input := LibraryRpc_SelectLabelsAll_Input {}
	_output := LibraryRpc_SelectLabelsAll_Output {}
	if _error := _client.rpc.Call ("Library.SelectLabelsAll", &_input, &_output); _error == nil {
		return _output.Labels, _output.Error
	} else {
		return nil, errorw (0xa8d5d0c5, _error)
	}
}

func (_exports *LibraryRpcServerExports) SelectLabelsAll (_input *LibraryRpc_SelectLabelsAll_Input, _output *LibraryRpc_SelectLabelsAll_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Labels, _output.Error = _exports.server.library.SelectLabelsAll ()
	return nil
}




type LibraryRpc_ResolveFullByLabel_Input struct {
	Label string
}

type LibraryRpc_ResolveFullByLabel_Output struct {
	Scriptlet *Scriptlet
	Error *Error
}

func (_client *LibraryRpcClient) ResolveFullByLabel (_label string) (*Scriptlet, *Error) {
	_input := LibraryRpc_ResolveFullByLabel_Input {}
	_output := LibraryRpc_ResolveFullByLabel_Output {}
	_input.Label = _label
	if _error := _client.rpc.Call ("Library.ResolveFullByLabel", &_input, &_output); _error == nil {
		return _output.Scriptlet, _output.Error
	} else {
		return nil, errorw (0x80b71741, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveFullByLabel (_input *LibraryRpc_ResolveFullByLabel_Input, _output *LibraryRpc_ResolveFullByLabel_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Scriptlet, _output.Error = _exports.server.library.ResolveFullByLabel (_input.Label)
	return nil
}




type LibraryRpc_ResolveMetaByLabel_Input struct {
	Label string
}

type LibraryRpc_ResolveMetaByLabel_Output struct {
	Scriptlet *Scriptlet
	Error *Error
}

func (_client *LibraryRpcClient) ResolveMetaByLabel (_label string) (*Scriptlet, *Error) {
	_input := LibraryRpc_ResolveMetaByLabel_Input {}
	_output := LibraryRpc_ResolveMetaByLabel_Output {}
	_input.Label = _label
	if _error := _client.rpc.Call ("Library.ResolveMetaByLabel", &_input, &_output); _error == nil {
		return _output.Scriptlet, _output.Error
	} else {
		return nil, errorw (0x7133a4fd, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveMetaByLabel (_input *LibraryRpc_ResolveMetaByLabel_Input, _output *LibraryRpc_ResolveMetaByLabel_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Scriptlet, _output.Error = _exports.server.library.ResolveMetaByLabel (_input.Label)
	return nil
}




type LibraryRpc_ResolveBodyByLabel_Input struct {
	Label string
}

type LibraryRpc_ResolveBodyByLabel_Output struct {
	Body string
	Found bool
	Error *Error
}

func (_client *LibraryRpcClient) ResolveBodyByLabel (_label string) (string, bool, *Error) {
	_input := LibraryRpc_ResolveBodyByLabel_Input {}
	_output := LibraryRpc_ResolveBodyByLabel_Output {}
	_input.Label = _label
	if _error := _client.rpc.Call ("Library.ResolveBodyByLabel", &_input, &_output); _error == nil {
		return _output.Body, _output.Found, _output.Error
	} else {
		return "", false, errorw (0x5d8357ba, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveBodyByLabel (_input *LibraryRpc_ResolveBodyByLabel_Input, _output *LibraryRpc_ResolveBodyByLabel_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Body, _output.Found, _output.Error = _exports.server.library.ResolveBodyByLabel (_input.Label)
	return nil
}




type LibraryRpc_ResolveFingerprintByLabel_Input struct {
	Label string
}

type LibraryRpc_ResolveFingerprintByLabel_Output struct {
	Fingerprint string
	Found bool
	Error *Error
}

func (_client *LibraryRpcClient) ResolveFingerprintByLabel (_label string) (string, bool, *Error) {
	_input := LibraryRpc_ResolveFingerprintByLabel_Input {}
	_output := LibraryRpc_ResolveFingerprintByLabel_Output {}
	_input.Label = _label
	if _error := _client.rpc.Call ("Library.ResolveFingerprintByLabel", &_input, &_output); _error == nil {
		return _output.Fingerprint, _output.Found, _output.Error
	} else {
		return "", false, errorw (0xc4df3289, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveFingerprintByLabel (_input *LibraryRpc_ResolveFingerprintByLabel_Input, _output *LibraryRpc_ResolveFingerprintByLabel_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Fingerprint, _output.Found, _output.Error = _exports.server.library.ResolveFingerprintByLabel (_input.Label)
	return nil
}



type LibraryRpc_ResolveContextByIdentifier_Input struct {
	Fingerprint string
}

type LibraryRpc_ResolveContextByIdentifier_Output struct {
	Context *ScriptletContext
	Found bool
	Error *Error
}

func (_client *LibraryRpcClient) ResolveContextByIdentifier (_fingerprint string) (*ScriptletContext, bool, *Error) {
	_input := LibraryRpc_ResolveContextByIdentifier_Input {}
	_output := LibraryRpc_ResolveContextByIdentifier_Output {}
	_input.Fingerprint = _fingerprint
	if _error := _client.rpc.Call ("Library.ResolveContextByIdentifier", &_input, &_output); _error == nil {
		return _output.Context, _output.Found, _output.Error
	} else {
		return nil, false, errorw (0x532bfaea, _error)
	}
}

func (_exports *LibraryRpcServerExports) ResolveContextByIdentifier (_input *LibraryRpc_ResolveContextByIdentifier_Input, _output *LibraryRpc_ResolveContextByIdentifier_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Context, _output.Found, _output.Error = _exports.server.library.ResolveContextByIdentifier (_input.Fingerprint)
	return nil
}




type LibraryRpc_SelectLibrarySources_Input struct {
}

type LibraryRpc_SelectLibrarySources_Output struct {
	Sources LibrarySources
	Error *Error
}

func (_client *LibraryRpcClient) SelectLibrarySources () (LibrarySources, *Error) {
	_input := LibraryRpc_SelectLibrarySources_Input {}
	_output := LibraryRpc_SelectLibrarySources_Output {}
	if _error := _client.rpc.Call ("Library.SelectLibrarySources", &_input, &_output); _error == nil {
		return _output.Sources, _output.Error
	} else {
		return nil, errorw (0x607c2934, _error)
	}
}

func (_exports *LibraryRpcServerExports) SelectLibrarySources (_input *LibraryRpc_SelectLibrarySources_Input, _output *LibraryRpc_SelectLibrarySources_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Sources, _output.Error = _exports.server.library.SelectLibrarySources ()
	return nil
}




type LibraryRpc_SelectLibraryContext_Input struct {
}

type LibraryRpc_SelectLibraryContext_Output struct {
	Context *LibraryContext
	Error *Error
}

func (_client *LibraryRpcClient) SelectLibraryContext () (*LibraryContext, *Error) {
	_input := LibraryRpc_SelectLibraryContext_Input {}
	_output := LibraryRpc_SelectLibraryContext_Output {}
	if _error := _client.rpc.Call ("Library.SelectLibraryContext", &_input, &_output); _error == nil {
		return _output.Context, _output.Error
	} else {
		return nil, errorw (0x4ec2bf5a, _error)
	}
}

func (_exports *LibraryRpcServerExports) SelectLibraryContext (_input *LibraryRpc_SelectLibraryContext_Input, _output *LibraryRpc_SelectLibraryContext_Output) (error) {
	_exports.server.mutex.RLock ()
	defer _exports.server.mutex.RUnlock ()
	_output.Context, _output.Error = _exports.server.library.SelectLibraryContext ()
	return nil
}




func urlSplit (_url string) (string, string, *Error) {
	_urlParts := strings.SplitAfterN (_url, ":", 2)
	if (len (_urlParts) < 1) || (len (_urlParts[0]) <= 1) || (len (_urlParts[1]) == 0) {
		return "", "", errorf (0x2dbbd8c7, "invalid URL")
	}
	return _urlParts[0][: len (_urlParts[0]) - 1], _urlParts[1], nil
}

