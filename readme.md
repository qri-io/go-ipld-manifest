### Manifests

For every discrete DAG _g_ one can construct a _manifest_ which is a second DAG of only block names and links (no content):

![graph_manifest](https://user-images.githubusercontent.com/1154390/45841799-e11cdd80-bce8-11e8-9fd9-572844aabf59.png)

These manifests are relatively small. If expressed as a set of two lists (one of array-positional links and one of names/hashes) it should be possible to represent many gigs worth of IPFS DAG content in < 100kb of CBOR.

IMHO, the power of IPFS is derived from the dual expression of blocks as both graphs and flat lists. This is also a fault line that shows up in the seam between bitswap and graph sync. I think graph manifests are a missing "primitive" from IPFS.

These manifests have a few properties that are nice:
* deterministic: a properly designed algo for generating manifests will generate the same manifest when given the same graph. Hash it, pass it to your friends. If the graph you're generating a manifest is immutable, manifest and hash of manifest are also immutable
* one can generate a manifest of any subgraph
* if implemented as a protocol, manifest generation can act as a challenge. You don't know if I've already generated this manifest ahead of time, and am simply asking you to compute it for trust purposes
* When I _don't_ have a manifest, I can ask multiple peers for the manifest of the same graph root. Differing responses raise suspicion
* _sooooooo cachable_, manifests could themselves be designed/required to fit into a single IPFS block.

If I'm planning on efficiently planning my requests for blocks, I _really_ want this manifest as soon as possible. Once I have a manifest I can trust I know a shit tonne of important things:
* when I'm done (progress)
* what blocks to ask for
* how blocks are related

So this might be a graph-sync thing, but it could also be a structural outgrowth of a bitswap session: establish a trusted graph, then divy up block requests among the session. If block sizes are also in the manifest, one can match larger blocks to faster peers. The point being, a manifest gives me a primitive to plan my block requests, and makes optimizing request planning a matter of better matching 

Downsides:
* you need the entire graph to calculate a manifest, or at least a trusted list of names and links (you may be able to use manifests to generate other manifests... a story for another day)
* graphs aren't super trivial to calculate, I could trick others into doing work they don't want to do if not rate limited or something. it's worth noting that calculating a manifest should be as cheap or cheaper as than through the block graph (cheaper if I can avoid loading associated data).

Both of those downsides can be mitigated by implementing manifests as a protocol, where peers can dynamically generate manifests of arbitrary graphs & subgraphs, which is the only reason I think it should exist at the IPFS layer.

Adding in Graph manifests is kinda like turning IPFS into dynamic bittorrent 🤷‍♂️.


### Manifest Sizes:
this repo contains a strawman test, runnable with:
```shell
go test -v
```
output:
```
manifest representing 4043 nodes and 1.024210Gb of content is 253.921997kb as CBOR
```

I need to double check this test for correctness (node count seems off to me), but it's a start.