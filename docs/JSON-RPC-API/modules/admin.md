






| Entity | Version |
| --- | --- |
| Source | <code>1.13.0-stable/generated-at:2026-09-14T07:51:25-06:00</code> |
| OpenRPC | <code>1.2.6</code> |

---




### admin_addPeer

AddPeer requests connecting to a remote node, and also maintaining the new
connection at all times, even reconnecting if it is lost.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_addPeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_addPeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.addPeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) AddPeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.AddPeer(node)
	return true, nil
}// AddPeer requests connecting to a remote node, and also maintaining the new
// connection at all times, even reconnecting if it is lost.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L61" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_addTrustedPeer

AddTrustedPeer allows a remote node to always connect, even if slots are full


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_addTrustedPeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_addTrustedPeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.addTrustedPeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) AddTrustedPeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.AddTrustedPeer(node)
	return true, nil
}// AddTrustedPeer allows a remote node to always connect, even if slots are full

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L93" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_datadir

Datadir retrieves the current data directory the node is using.


#### Params (0)

_None_

#### Result




<code>string</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_datadir", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_datadir", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.datadir();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) Datadir() string {
	return api.node.DataDir()
}// Datadir retrieves the current data directory the node is using.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L324" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_ecbp1100

Ecbp1100 is MessActivate under its ECBP-1100 name, which it had first.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `The block height description`
			- enum: earliest, latest, pending
			- title: `blockNumberTag`
			- type: string


			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: string


	- title: `blockNumberIdentifier`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "The block height description",
                "enum": [
                    "earliest",
                    "latest",
                    "pending"
                ],
                "title": "blockNumberTag",
                "type": [
                    "string"
                ]
            },
            {
                "description": "Hex representation of a uint64",
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": [
                    "string"
                ]
            }
        ],
        "title": "blockNumberIdentifier"
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_ecbp1100", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_ecbp1100", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	admin.ecbp1100(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) Ecbp1100(blockNr rpc.BlockNumber) (bool, error) {
	return api.MessActivate(blockNr)
}// Ecbp1100 is MessActivate under its ECBP-1100 name, which it had first.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L437" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_ecbp1100Deactivate

Ecbp1100Deactivate is MessDeactivate under its ECBP-1100 name.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `The block height description`
			- enum: earliest, latest, pending
			- title: `blockNumberTag`
			- type: string


			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: string


	- title: `blockNumberIdentifier`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "The block height description",
                "enum": [
                    "earliest",
                    "latest",
                    "pending"
                ],
                "title": "blockNumberTag",
                "type": [
                    "string"
                ]
            },
            {
                "description": "Hex representation of a uint64",
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": [
                    "string"
                ]
            }
        ],
        "title": "blockNumberIdentifier"
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_ecbp1100Deactivate", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_ecbp1100Deactivate", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	admin.ecbp1100Deactivate(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) Ecbp1100Deactivate(blockNr rpc.BlockNumber) (bool, error) {
	return api.MessDeactivate(blockNr)
}// Ecbp1100Deactivate is MessDeactivate under its ECBP-1100 name.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L442" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_ecbp1100NoDisable

Ecbp1100NoDisable is MessNoDisable under its ECBP-1100 name.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
enable <code>bool</code> 

  + Required: ✓ Yes






#### Result




<code>ECBP1100Status</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- activatedAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- defaultDisabledAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- enabled: 
			- type: `boolean`

		- head: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- noDisable: 
			- type: `boolean`

		- nodeSwitch: 
			- type: `boolean`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "activatedAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "defaultDisabledAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "enabled": {
                "type": "boolean"
            },
            "head": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "noDisable": {
                "type": "boolean"
            },
            "nodeSwitch": {
                "type": "boolean"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_ecbp1100NoDisable", "params": [<enable>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_ecbp1100NoDisable", "params": [<enable>]}'
	```


=== "Javascript Console"

	``` js
	admin.ecbp1100NoDisable(enable);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) Ecbp1100NoDisable(enable bool) (ECBP1100Status, error) {
	return api.MessNoDisable(enable)
}// Ecbp1100NoDisable is MessNoDisable under its ECBP-1100 name.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L447" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_ecbp1100Status

Ecbp1100Status is MessStatus under its ECBP-1100 name, which it had first.


#### Params (0)

_None_

#### Result




<code>ECBP1100Status</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- activatedAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- defaultDisabledAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- enabled: 
			- type: `boolean`

		- head: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- noDisable: 
			- type: `boolean`

		- nodeSwitch: 
			- type: `boolean`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "activatedAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "defaultDisabledAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "enabled": {
                "type": "boolean"
            },
            "head": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "noDisable": {
                "type": "boolean"
            },
            "nodeSwitch": {
                "type": "boolean"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_ecbp1100Status", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_ecbp1100Status", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.ecbp1100Status();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) Ecbp1100Status() ECBP1100Status {
	return api.MessStatus()
}// Ecbp1100Status is MessStatus under its ECBP-1100 name, which it had first.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L452" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_exportChain

ExportChain exports the current blockchain into a local file,
or a range of blocks if first and last are non-nil.


#### Params (3)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes





__2:__ 
first <code>*uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
last <code>*uint64</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_exportChain", "params": [<file>, <first>, <last>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_exportChain", "params": [<file>, <first>, <last>]}'
	```


=== "Javascript Console"

	``` js
	admin.exportChain(file,first,last);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) ExportChain(file string, first *uint64, last *uint64) (bool, error) {
	if first == nil && last != nil {
		return false, errors.New("last cannot be specified without first")
	}
	if first != nil && last == nil {
		head := api.eth.BlockChain().CurrentHeader().Number.Uint64()
		last = &head
	}
	if _, err := os.Stat(file); err == nil {
		return false, errors.New("location would overwrite an existing file")
	}
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false, err
	}
	defer out.Close()
	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}
	if first != nil {
		if err := api.eth.BlockChain().ExportN(writer, *first, *last); err != nil {
			return false, err
		}
	} else if err := api.eth.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}// ExportChain exports the current blockchain into a local file,
// or a range of blocks if first and last are non-nil.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L49" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_importChain

ImportChain imports a blockchain from a local file.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
file <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_importChain", "params": [<file>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_importChain", "params": [<file>]}'
	```


=== "Javascript Console"

	``` js
	admin.importChain(file);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) ImportChain(file string) (bool, error) {
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()
	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}
	stream := rlp.NewStream(reader, 0)
	blocks, index := make([ // ImportChain imports a blockchain from a local file.
	]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			if block.NumberU64() == 0 {
				continue
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}
		if hasAllBlocks(api.eth.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}
```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L97" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_maxPeers

MaxPeers sets the maximum peer limit for the protocol manager and the p2p server.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
n <code>int</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_maxPeers", "params": [<n>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_maxPeers", "params": [<n>]}'
	```


=== "Javascript Console"

	``` js
	admin.maxPeers(n);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) MaxPeers(n int) (bool, error) {
	api.eth.handler.maxPeers = n
	api.eth.p2pServer.MaxPeers = n
	for i := api.eth.handler.peers.len(); i > n; i = api.eth.handler.peers.len() {
		p := api.eth.handler.peers.WorstPeer()
		if p == nil {
			break
		}
		api.eth.handler.removePeer(p.ID())
	}
	return true, nil
}// MaxPeers sets the maximum peer limit for the protocol manager and the p2p server.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L238" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_mess

Mess turns the ECBP-1100 (MESS) chain-selection defense on or off in the running node and
reports the state afterwards. It is the runtime counterpart of --mess and --mess=false, and
it works whatever the two blocks are, so the caller does not have to know which one keeps
MESS off.

On pushes the deactivation out of reach and, if the activation is unset or sits beyond the
head, as --mess=false leaves it, moves it back to the network's own activation block, as
--mess does, or to block 0 where the network ships none or the head has not reached it. Off
pushes the activation out of reach, so MESS does not apply whatever the deactivation block is.

It does not touch the node-level switch that the low-peer-count and stale-head safeguards
turn off, so Enabled can still be false with the window open. NodeSwitch reports that switch,
and MessNoDisable or --mess.nodisable keeps it on.

This mutates chain configuration, and does not persist: the flags and the config file
decide again at the next start.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
enable <code>bool</code> 

  + Required: ✓ Yes






#### Result




<code>ECBP1100Status</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- activatedAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- defaultDisabledAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- enabled: 
			- type: `boolean`

		- head: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- noDisable: 
			- type: `boolean`

		- nodeSwitch: 
			- type: `boolean`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "activatedAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "defaultDisabledAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "enabled": {
                "type": "boolean"
            },
            "head": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "noDisable": {
                "type": "boolean"
            },
            "nodeSwitch": {
                "type": "boolean"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_mess", "params": [<enable>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_mess", "params": [<enable>]}'
	```


=== "Javascript Console"

	``` js
	admin.mess(enable);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) Mess(enable bool) (ECBP1100Status, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	if err := applyMESSSwitch(api.eth.blockchain.Config(), head.Uint64(), api.eth.messBundledActivation, enable); err != nil {
		return ECBP1100Status{}, err
	}
	return api.MessStatus(), nil
}// Mess turns the ECBP-1100 (MESS) chain-selection defense on or off in the running node and
// reports the state afterwards. It is the runtime counterpart of --mess and --mess=false, and
// it works whatever the two blocks are, so the caller does not have to know which one keeps
// MESS off.
//
// On pushes the deactivation out of reach and, if the activation is unset or sits beyond the
// head, as --mess=false leaves it, moves it back to the network's own activation block, as
// --mess does, or to block 0 where the network ships none or the head has not reached it. Off
// pushes the activation out of reach, so MESS does not apply whatever the deactivation block is.
//
// It does not touch the node-level switch that the low-peer-count and stale-head safeguards
// turn off, so Enabled can still be false with the window open. NodeSwitch reports that switch,
// and MessNoDisable or --mess.nodisable keeps it on.
//
// This mutates chain configuration, and does not persist: the flags and the config file
// decide again at the next start.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L257" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_messActivate

MessActivate sets the ECBP-1100 (MESS) activation block and reports whether MESS applies
afterwards. The block is a height, or "latest" or "pending", which both mean the current head;
"finalized" and "safe" are refused.

It is the runtime counterpart of --mess.activate, and does what the flag does: MESS applies
from the block with no end, unless a deactivation block was given, with --mess.deactivate, a
config file line or MessDeactivate. Then the later of the two blocks decides: an activation
after the deactivation block turns MESS on, and one at or before it leaves MESS off once the
head reaches the deactivation block, which the node then logs. Mess turns MESS on or off
whatever the two blocks are.

This mutates chain configuration, and does not persist. To read the current state without
changing it, use MessStatus.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `The block height description`
			- enum: earliest, latest, pending
			- title: `blockNumberTag`
			- type: string


			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: string


	- title: `blockNumberIdentifier`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "The block height description",
                "enum": [
                    "earliest",
                    "latest",
                    "pending"
                ],
                "title": "blockNumberTag",
                "type": [
                    "string"
                ]
            },
            {
                "description": "Hex representation of a uint64",
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": [
                    "string"
                ]
            }
        ],
        "title": "blockNumberIdentifier"
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_messActivate", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_messActivate", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	admin.messActivate(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) MessActivate(blockNr rpc.BlockNumber) (bool, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	i, err := ecbp1100ActivationBlock(blockNr, head.Uint64())
	if err != nil {
		return false, err
	}
	if err := applyMESSActivation(api.eth.blockchain.Config(), i, api.messDeactivationGiven.Load()); err != nil {
		return false, err
	}
	status := api.MessStatus()
	if !status.Enabled && status.NodeSwitch {
		if d := api.eth.blockchain.Config().GetECBP1100DeactivateTransition(); messKeptOffBy(i, d, head.Uint64()) {
			log.Warn("ECBP1100 (MESS) stays off: the deactivation block is at or after the activation block set, and the chain has reached it", "activation", i, "deactivation", *d, "head", head, "hint", "admin_mess(true) turns MESS on")
		}
	}
	return status.Enabled, nil
}// MessActivate sets the ECBP-1100 (MESS) activation block and reports whether MESS applies
// afterwards. The block is a height, or "latest" or "pending", which both mean the current head;
// "finalized" and "safe" are refused.
//
// It is the runtime counterpart of --mess.activate, and does what the flag does: MESS applies
// from the block with no end, unless a deactivation block was given, with --mess.deactivate, a
// config file line or MessDeactivate. Then the later of the two blocks decides: an activation
// after the deactivation block turns MESS on, and one at or before it leaves MESS off once the
// head reaches the deactivation block, which the node then logs. Mess turns MESS on or off
// whatever the two blocks are.
//
// This mutates chain configuration, and does not persist. To read the current state without
// changing it, use MessStatus.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L172" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_messDeactivate

MessDeactivate sets the ECBP-1100 (MESS) deactivation block and reports whether MESS applies
afterwards. The block is read as for MessActivate.

It is the runtime counterpart of --mess.deactivate. Of the activation and the deactivation
block, the later one decides: a deactivation at or after the activation turns MESS off from
that block, and one before the activation does not turn it off. MessActivate and MessNoDisable
keep a deactivation set here, as --mess.activate keeps one given with --mess.deactivate.

This mutates chain configuration, and does not persist. To read the current state without
changing it, use MessStatus.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
blockNr <code>rpc.BlockNumber</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- oneOf: 

			- description: `The block height description`
			- enum: earliest, latest, pending
			- title: `blockNumberTag`
			- type: string


			- description: `Hex representation of a uint64`
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: string


	- title: `blockNumberIdentifier`


	```

=== "Raw"

	``` Raw
	{
        "oneOf": [
            {
                "description": "The block height description",
                "enum": [
                    "earliest",
                    "latest",
                    "pending"
                ],
                "title": "blockNumberTag",
                "type": [
                    "string"
                ]
            },
            {
                "description": "Hex representation of a uint64",
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": [
                    "string"
                ]
            }
        ],
        "title": "blockNumberIdentifier"
    }
	```





#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_messDeactivate", "params": [<blockNr>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_messDeactivate", "params": [<blockNr>]}'
	```


=== "Javascript Console"

	``` js
	admin.messDeactivate(blockNr);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) MessDeactivate(blockNr rpc.BlockNumber) (bool, error) {
	head := api.eth.blockchain.CurrentBlock().Number
	i, err := ecbp1100DeactivationBlock(blockNr, head.Uint64())
	if err != nil {
		return false, err
	}
	if err := api.eth.blockchain.Config().SetECBP1100DeactivateTransition(&i); err != nil {
		return false, err
	}
	api.messDeactivationGiven.Store(true)
	return api.MessStatus().Enabled, nil
}// MessDeactivate sets the ECBP-1100 (MESS) deactivation block and reports whether MESS applies
// afterwards. The block is read as for MessActivate.
//
// It is the runtime counterpart of --mess.deactivate. Of the activation and the deactivation
// block, the later one decides: a deactivation at or after the activation turns MESS off from
// that block, and one before the activation does not turn it off. MessActivate and MessNoDisable
// keep a deactivation set here, as --mess.activate keeps one given with --mess.deactivate.
//
// This mutates chain configuration, and does not persist. To read the current state without
// changing it, use MessStatus.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L228" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_messNoDisable

MessNoDisable is the runtime counterpart of --mess.nodisable, and reports the state afterwards.

True does what the flag does. It turns MESS on: an activation that keeps MESS off moves as Mess
moves it, and the deactivation moves out of reach unless one was given. And it keeps MESS on
once it is on: the low-peer-count and stale-head safeguards no longer switch it off. Like the
flag, it does not switch MESS on while those safeguards hold it off; MESS comes on at the next
sync, and then stays on.

False lets the safeguards switch MESS off again, as --mess.nodisable=false does, and leaves the
two blocks alone.

This mutates chain configuration, and does not persist: the flags and the config file
decide again at the next start.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
enable <code>bool</code> 

  + Required: ✓ Yes






#### Result




<code>ECBP1100Status</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- activatedAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- defaultDisabledAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- enabled: 
			- type: `boolean`

		- head: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- noDisable: 
			- type: `boolean`

		- nodeSwitch: 
			- type: `boolean`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "activatedAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "defaultDisabledAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "enabled": {
                "type": "boolean"
            },
            "head": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "noDisable": {
                "type": "boolean"
            },
            "nodeSwitch": {
                "type": "boolean"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_messNoDisable", "params": [<enable>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_messNoDisable", "params": [<enable>]}'
	```


=== "Javascript Console"

	``` js
	admin.messNoDisable(enable);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) MessNoDisable(enable bool) (ECBP1100Status, error) {
	if !enable {
		api.eth.blockchain.ArtificialFinalityNoDisable(0)
		return api.MessStatus(), nil
	}
	head := api.eth.blockchain.CurrentBlock().Number.Uint64()
	if err := applyMESSNoDisable(api.eth.blockchain.Config(), head, api.eth.messBundledActivation, api.messDeactivationGiven.Load()); err != nil {
		return ECBP1100Status{}, err
	}
	api.eth.blockchain.ArtificialFinalityNoDisable(1)
	return api.MessStatus(), nil
}// MessNoDisable is the runtime counterpart of --mess.nodisable, and reports the state afterwards.
//
// True does what the flag does. It turns MESS on: an activation that keeps MESS off moves as Mess
// moves it, and the deactivation moves out of reach unless one was given. And it keeps MESS on
// once it is on: the low-peer-count and stale-head safeguards no longer switch it off. Like the
// flag, it does not switch MESS on while those safeguards hold it off; MESS comes on at the next
// sync, and then stays on.
//
// False lets the safeguards switch MESS off again, as --mess.nodisable=false does, and leaves the
// two blocks alone.
//
// This mutates chain configuration, and does not persist: the flags and the config file
// decide again at the next start.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L316" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_messStatus

MessStatus reports the ECBP-1100 (MESS) state at the current head and changes nothing.

MessActivate answers a similar question by first assigning the activation block it is
passed, so it cannot be used to observe a running node. This exists so that state can be read
without altering it.


#### Params (0)

_None_

#### Result




<code>ECBP1100Status</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- activatedAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- defaultDisabledAtBlock: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- enabled: 
			- type: `boolean`

		- head: 
			- pattern: `^0x([a-fA-F\d])+$`
			- title: `uint64`
			- type: `string`

		- noDisable: 
			- type: `boolean`

		- nodeSwitch: 
			- type: `boolean`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "activatedAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "defaultDisabledAtBlock": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "enabled": {
                "type": "boolean"
            },
            "head": {
                "pattern": "^0x([a-fA-F\\d])+$",
                "title": "uint64",
                "type": "string"
            },
            "noDisable": {
                "type": "boolean"
            },
            "nodeSwitch": {
                "type": "boolean"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_messStatus", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_messStatus", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.messStatus();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *AdminAPI) MessStatus() ECBP1100Status {
	bc := api.eth.blockchain
	return ecbp1100Status(bc.Config(), bc.CurrentBlock().Number, bc.IsArtificialFinalityEnabled(), bc.IsArtificialFinalityNoDisable())
}// MessStatus reports the ECBP-1100 (MESS) state at the current head and changes nothing.
//
// MessActivate answers a similar question by first assigning the activation block it is
// passed, so it cannot be used to observe a running node. This exists so that state can be read
// without altering it.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/eth/api_admin.go#L415" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_nodeInfo

NodeInfo retrieves all the information we know about the host node at the
protocol granularity.


#### Params (0)

_None_

#### Result




<code>*p2p.NodeInfo</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- additionalProperties: `false`
	- properties: 
		- enode: 
			- type: `string`

		- enr: 
			- type: `string`

		- id: 
			- type: `string`

		- ip: 
			- type: `string`

		- listenAddr: 
			- type: `string`

		- name: 
			- type: `string`

		- ports: 
			- additionalProperties: `false`
			- properties: 
				- discovery: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- listener: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`


			- type: `object`

		- protocols: 
			- additionalProperties: `false`
			- properties: 
				- discovery: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`

				- listener: 
					- pattern: `^0x[a-fA-F0-9]+$`
					- title: `integer`
					- type: `string`


			- type: `object`


	- type: object


	```

=== "Raw"

	``` Raw
	{
        "additionalProperties": false,
        "properties": {
            "enode": {
                "type": "string"
            },
            "enr": {
                "type": "string"
            },
            "id": {
                "type": "string"
            },
            "ip": {
                "type": "string"
            },
            "listenAddr": {
                "type": "string"
            },
            "name": {
                "type": "string"
            },
            "ports": {
                "additionalProperties": false,
                "properties": {
                    "discovery": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "listener": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    }
                },
                "type": "object"
            },
            "protocols": {
                "additionalProperties": false,
                "properties": {
                    "discovery": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    },
                    "listener": {
                        "pattern": "^0x[a-fA-F0-9]+$",
                        "title": "integer",
                        "type": "string"
                    }
                },
                "type": "object"
            }
        },
        "type": [
            "object"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_nodeInfo", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_nodeInfo", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.nodeInfo();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) NodeInfo() (*p2p.NodeInfo, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	return server.NodeInfo(), nil
}// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L314" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_peerEvents

PeerEvents creates an RPC subscription which receives peer events from the
node's p2p.Server


#### Params (0)

_None_

#### Result




<code>*rpc.Subscription</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Subscription identifier`
	- title: `subscriptionID`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Subscription identifier",
        "title": "subscriptionID",
        "type": [
            "string"
        ]
    }
	```



#### Client Method Invocation Examples








=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_subscribe", "params": ["peerEvents"]}'
	```




<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) PeerEvents(ctx context.Context) (*rpc.Subscription, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}
	rpcSub := notifier.CreateSubscription()
	go func() {
		events := make(chan *p2p.PeerEvent)
		sub := server.SubscribeEvents(events)
		defer sub.Unsubscribe()
		for {
			select {
			case event := <-events:
				notifier.Notify(rpcSub.ID, event)
			case <-sub.Err():
				return
			case <-rpcSub.Err():
				return
			case <-notifier.Closed():
				return
			}
		}
	}()
	return rpcSub, nil
}// PeerEvents creates an RPC subscription which receives peer events from the
// node's p2p.Server

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L125" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_peers

Peers retrieves all the information we know about each individual peer at the
protocol granularity.


#### Params (0)

_None_

#### Result



p2pPeerInfo <code>[]*p2p.PeerInfo</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- items: 

			- additionalProperties: `false`
			- properties: 
				- caps: 
					- items: 
						- type: `string`

					- type: `array`

				- enode: 
					- type: `string`

				- enr: 
					- type: `string`

				- id: 
					- type: `string`

				- name: 
					- type: `string`

				- network: 
					- additionalProperties: `false`
					- properties: 
						- inbound: 
							- type: `boolean`

						- localAddress: 
							- type: `string`

						- remoteAddress: 
							- type: `string`

						- static: 
							- type: `boolean`

						- trusted: 
							- type: `boolean`


					- type: `object`

				- protocols: 
					- additionalProperties: `false`
					- properties: 
						- inbound: 
							- type: `boolean`

						- localAddress: 
							- type: `string`

						- remoteAddress: 
							- type: `string`

						- static: 
							- type: `boolean`

						- trusted: 
							- type: `boolean`


					- type: `object`


			- type: object


	- type: array


	```

=== "Raw"

	``` Raw
	{
        "items": [
            {
                "additionalProperties": false,
                "properties": {
                    "caps": {
                        "items": {
                            "type": "string"
                        },
                        "type": "array"
                    },
                    "enode": {
                        "type": "string"
                    },
                    "enr": {
                        "type": "string"
                    },
                    "id": {
                        "type": "string"
                    },
                    "name": {
                        "type": "string"
                    },
                    "network": {
                        "additionalProperties": false,
                        "properties": {
                            "inbound": {
                                "type": "boolean"
                            },
                            "localAddress": {
                                "type": "string"
                            },
                            "remoteAddress": {
                                "type": "string"
                            },
                            "static": {
                                "type": "boolean"
                            },
                            "trusted": {
                                "type": "boolean"
                            }
                        },
                        "type": "object"
                    },
                    "protocols": {
                        "additionalProperties": false,
                        "properties": {
                            "inbound": {
                                "type": "boolean"
                            },
                            "localAddress": {
                                "type": "string"
                            },
                            "remoteAddress": {
                                "type": "string"
                            },
                            "static": {
                                "type": "boolean"
                            },
                            "trusted": {
                                "type": "boolean"
                            }
                        },
                        "type": "object"
                    }
                },
                "type": [
                    "object"
                ]
            }
        ],
        "type": [
            "array"
        ]
    }
	```



#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_peers", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_peers", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.peers();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) Peers() ([ // Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
]*p2p.PeerInfo, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	return server.PeersInfo(), nil
}
```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L304" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_removePeer

RemovePeer disconnects from a remote node if the connection exists


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_removePeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_removePeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.removePeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) RemovePeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.RemovePeer(node)
	return true, nil
}// RemovePeer disconnects from a remote node if the connection exists

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L77" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_removeTrustedPeer

RemoveTrustedPeer removes a remote node from the trusted peer set, but it
does not disconnect it automatically.


#### Params (1)

Parameters must be given _by position_.


__1:__ 
url <code>string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_removeTrustedPeer", "params": [<url>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_removeTrustedPeer", "params": [<url>]}'
	```


=== "Javascript Console"

	``` js
	admin.removeTrustedPeer(url);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) RemoveTrustedPeer(url string) (bool, error) {
	server := api.node.Server()
	if server == nil {
		return false, ErrNodeStopped
	}
	node, err := enode.Parse(enode.ValidSchemes, url)
	if err != nil {
		return false, fmt.Errorf("invalid enode: %v", err)
	}
	server.RemoveTrustedPeer(node)
	return true, nil
}// RemoveTrustedPeer removes a remote node from the trusted peer set, but it
// does not disconnect it automatically.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L109" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_startHTTP

StartHTTP starts the HTTP RPC API server.


#### Params (5)

Parameters must be given _by position_.


__1:__ 
host <code>*string</code> 

  + Required: ✓ Yes





__2:__ 
port <code>*int</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
cors <code>*string</code> 

  + Required: ✓ Yes





__4:__ 
apis <code>*string</code> 

  + Required: ✓ Yes





__5:__ 
vhosts <code>*string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_startHTTP", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_startHTTP", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```


=== "Javascript Console"

	``` js
	admin.startHTTP(host,port,cors,apis,vhosts);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StartHTTP(host *string, port *int, cors *string, apis *string, vhosts *string) (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()
	if host == nil {
		h := DefaultHTTPHost
		if api.node.config.HTTPHost != "" {
			h = api.node.config.HTTPHost
		}
		host = &h
	}
	if port == nil {
		port = &api.node.config.HTTPPort
	}
	config := httpConfig{CorsAllowedOrigins: api.node.config.HTTPCors, Vhosts: api.node.config.HTTPVirtualHosts, Modules: api.node.config.HTTPModules, rpcEndpointConfig: rpcEndpointConfig{batchItemLimit: api.node.config.BatchRequestLimit, batchResponseSizeLimit: api.node.config.BatchResponseMaxSize}}
	if cors != nil {
		config.CorsAllowedOrigins = nil
		for _, origin := // StartHTTP starts the HTTP RPC API server.
		range strings.Split(*cors, ",") {
			config.CorsAllowedOrigins = append(config.CorsAllowedOrigins, strings.TrimSpace(origin))
		}
	}
	if vhosts != nil {
		config.Vhosts = nil
		for _, vhost := range strings.Split(*host, ",") {
			config.Vhosts = append(config.Vhosts, strings.TrimSpace(vhost))
		}
	}
	if apis != nil {
		config.Modules = nil
		for _, m := range strings.Split(*apis, ",") {
			config.Modules = append(config.Modules, strings.TrimSpace(m))
		}
	}
	if err := api.node.http.setListenAddr(*host, *port); err != nil {
		return false, err
	}
	if err := api.node.http.enableRPC(api.node.rpcAPIs, config); err != nil {
		return false, err
	}
	if err := api.node.http.start(); err != nil {
		return false, err
	}
	return true, nil
}
```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L162" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_startRPC

StartRPC starts the HTTP RPC API server.
Deprecated: use StartHTTP instead.


#### Params (5)

Parameters must be given _by position_.


__1:__ 
host <code>*string</code> 

  + Required: ✓ Yes





__2:__ 
port <code>*int</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
cors <code>*string</code> 

  + Required: ✓ Yes





__4:__ 
apis <code>*string</code> 

  + Required: ✓ Yes





__5:__ 
vhosts <code>*string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_startRPC", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_startRPC", "params": [<host>, <port>, <cors>, <apis>, <vhosts>]}'
	```


=== "Javascript Console"

	``` js
	admin.startRPC(host,port,cors,apis,vhosts);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StartRPC(host *string, port *int, cors *string, apis *string, vhosts *string) (bool, error) {
	log.Warn("Deprecation warning", "method", "admin.StartRPC", "use-instead", "admin.StartHTTP")
	return api.StartHTTP(host, port, cors, apis, vhosts)
}// StartRPC starts the HTTP RPC API server.
// Deprecated: use StartHTTP instead.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L221" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_startWS

StartWS starts the websocket RPC API server.


#### Params (4)

Parameters must be given _by position_.


__1:__ 
host <code>*string</code> 

  + Required: ✓ Yes





__2:__ 
port <code>*int</code> 

  + Required: ✓ Yes


=== "Schema"

	``` Schema
	
	- description: `Hex representation of the integer`
	- pattern: `^0x[a-fA-F0-9]+$`
	- title: `integer`
	- type: string


	```

=== "Raw"

	``` Raw
	{
        "description": "Hex representation of the integer",
        "pattern": "^0x[a-fA-F0-9]+$",
        "title": "integer",
        "type": [
            "string"
        ]
    }
	```




__3:__ 
allowedOrigins <code>*string</code> 

  + Required: ✓ Yes





__4:__ 
apis <code>*string</code> 

  + Required: ✓ Yes






#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_startWS", "params": [<host>, <port>, <allowedOrigins>, <apis>]}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_startWS", "params": [<host>, <port>, <allowedOrigins>, <apis>]}'
	```


=== "Javascript Console"

	``` js
	admin.startWS(host,port,allowedOrigins,apis);
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StartWS(host *string, port *int, allowedOrigins *string, apis *string) (bool, error) {
	api.node.lock.Lock()
	defer api.node.lock.Unlock()
	if host == nil {
		h := DefaultWSHost
		if api.node.config.WSHost != "" {
			h = api.node.config.WSHost
		}
		host = &h
	}
	if port == nil {
		port = &api.node.config.WSPort
	}
	config := wsConfig{Modules: api.node.config.WSModules, Origins: api.node.config.WSOrigins, rpcEndpointConfig: rpcEndpointConfig{batchItemLimit: api.node.config.BatchRequestLimit, batchResponseSizeLimit: api.node.config.BatchResponseMaxSize}}
	if apis != nil {
		config.Modules = nil
		for _, m := // StartWS starts the websocket RPC API server.
		range strings.Split(*apis, ",") {
			config.Modules = append(config.Modules, strings.TrimSpace(m))
		}
	}
	if allowedOrigins != nil {
		config.Origins = nil
		for _, origin := range strings.Split(*allowedOrigins, ",") {
			config.Origins = append(config.Origins, strings.TrimSpace(origin))
		}
	}
	server := api.node.wsServerForPort(*port, false)
	if err := server.setListenAddr(*host, *port); err != nil {
		return false, err
	}
	openApis, _ := api.node.getAPIs()
	if err := server.enableWS(openApis, config); err != nil {
		return false, err
	}
	if err := server.start(); err != nil {
		return false, err
	}
	api.node.http.log.Info("WebSocket endpoint opened", "url", api.node.WSEndpoint())
	return true, nil
}
```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L240" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_stopHTTP

StopHTTP shuts down the HTTP server.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_stopHTTP", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_stopHTTP", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.stopHTTP();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StopHTTP() (bool, error) {
	api.node.http.stop()
	return true, nil
}// StopHTTP shuts down the HTTP server.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L227" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_stopRPC

StopRPC shuts down the HTTP server.
Deprecated: use StopHTTP instead.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_stopRPC", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_stopRPC", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.stopRPC();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StopRPC() (bool, error) {
	log.Warn("Deprecation warning", "method", "admin.StopRPC", "use-instead", "admin.StopHTTP")
	return api.StopHTTP()
}// StopRPC shuts down the HTTP server.
// Deprecated: use StopHTTP instead.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L234" target="_">View on GitHub →</a>
</p>
</details>

---



### admin_stopWS

StopWS terminates all WebSocket servers.


#### Params (0)

_None_

#### Result




<code>bool</code> 

  + Required: ✓ Yes




#### Client Method Invocation Examples


=== "Shell HTTP"

	``` shell
	curl -X POST -H "Content-Type: application/json" http://localhost:8545 --data '{"jsonrpc": "2.0", "id": 42, "method": "admin_stopWS", "params": []}'
	```





=== "Shell WebSocket"

	``` shell
	wscat -c ws://localhost:8546 -x '{"jsonrpc": "2.0", "id": 1, "method": "admin_stopWS", "params": []}'
	```


=== "Javascript Console"

	``` js
	admin.stopWS();
	```



<details><summary>Source code</summary>
<p>
```go
func (api *adminAPI) StopWS() (bool, error) {
	api.node.http.stopWS()
	api.node.ws.stop()
	return true, nil
}// StopWS terminates all WebSocket servers.

```
<a href="https://github.com/ethereumclassic/core-geth/blob/main/node/api.go#L296" target="_">View on GitHub →</a>
</p>
</details>

---

