**This is a very preliminary draft!**

Proposal
========

Aereum should be a protocol capable of providing stable and resilient
descentralized infrastructure for publishing and dissemination of original 
content for controled audiences. The protocol is based on the following 
hypothesis:

* Advertising should be capable of providing real economic resources to
keep the necessary infrastructure of the network, remuneration for content 
creators and organic growth.

* All members, in all capacities, should be genuine stakeholders of the network.

* There should be no possibility of top-down censoring, but there should be
effective blacklisting mechanisms. 

* No hidden illegal business entirely based on the network should be possible.
The advertising content should be entirely public.

Aero
====

Aero is the name of the token of the protocol. It is created out of thin air, 
through a mining mechanism, provided either by a proof of work or proof of 
stake blockchain formation mechanism. 

It should be designed exclusively to be a means of exchange within the network,
through which real economic value is traded between actors, but should not be
relied as a long-term medium for storage of value.

In principle unlimited units of aero can be created as the monetary base grows
along with activity in network itself. Parameters should be calibrated in order
to provide a slightly inflationary behavior of aero over longer periods of time.

Nodes
=====

Nodes with direct participation on the network should distribute information
about recent blocks under appropriate conditions. 

Mining
======

Miners are responsible for blockchain formation, unequivocally aggregating into 
the chain database all the valid messages submitted to the network. They collect 
the minted coins created on block formation and clearing fees for the 
processed messages. Details about mining protocols will be established in 
appropriate time.

Power of attorney and Intermediaties
====================================

Since the purpose of the network is communication, if it is succesful, it must
be capable of processing hundreds of millions or even billions of messages per
day. This requiremente poses relevant burden on the infrastructure that must be
aleviated by requirent as minimum blockchain memory as possible. We thus expect 
the competitive fully distributed behavior of nodes to keep track only of 
recent activity and general information like wallet content, valid subscribers,
valid audiences, etc. 

The relevant history of the blockchain should be an economic activity of 
intermediaires that will offer services to end users. These services can
incorporate things like wallet and crypto management, timeline indexation, 
timeline algorithms, and so on. This should be pursued on free contractual terms
off-network. To facilitate the commercial link between users and
intermediaries we invoke the principle of power of attoney by which o token 
gives right to another token to sign messages in its behalf. The intermediare 
can enter on a revenue sharing deal with subscribers to share aero due to the
later. 

On the other hand, as long as the subscriber is in possession of the original
cryptographic keys he has the autonomy to cancel within the network the power
of attorney at his free will, notwithstanding contractual obligations 
established off-network under virtual or real jurisdictions.

Block Chain
===========

Blockchain provides distributed timestamping (by epoch of block formation) and
unique formation of messages databases.

Specification for the format, rules and validation for blockchain formation
will be provided.

Advertising
===========

Advertising is the link between the real economy and the digital network. Within
Aereum advertising is an act by which a content creator redistributes to its 
audience content created by the advertiser and recieved in return Aero tokens.
Part of the tokens collected from the advertiser is also distributed 
(uniformly?) to the members of the audience. This resource will provide
digital resources for the less active members of the community to pay for the 
associated costs. 

Messages
========

There are only two kind of messages valid for the block chain.

The first one, is only valid for the genesis block and defines

```
Genesis:
   Token
   Wallet Token
   IP Address 
   Epoch=0
```

Token becomes an acceptable author. Wallet token receives 100 Aero. IP Address 
is a valid IP Address to a functioning node for the network.

All subsequent messages should be eith an Aero transfer

```
Transfer:
    From Wallet Token
    To Wallet Token
    Value
    Epoch 
    Signature by From Wallet Token (Private Key) of the above
```

or a general message 

```
Message:
   UUID
   Author Token
   Message-Type
   Message-Data
   Publishing Fee:
      Publishing Fee UUID
      Wallet Token
      Value
      Signature by Wallet Token (Private Key) of the above
    Epoch
    Power of attorney Token
    Signature either by Author Token (Private Key) of the above or
    by the valid power of attorney token
    Signature by Wallet Token (Private Key) of the above
```

Author Token should be a token associated to a valid subscriber (see bellow).
A message should not be incorporated into the block chain if it is old enough 
(proposed epoch less than a certain number of epochs prior of the newly formed
blockchain), the any signature is invalid, Message type is invalid, there is not
sufficient resources in the wallet or, finally, Message-Data suprpases a 
certain limit. Power of attorney token should refer to a valid non expired and
not revoked power of attorney message in the block chain (see bellow).

A new block cannot contain repeated messages. 

Message Types
=============

A subscription request

```
Subscribe:
   Caption
   Subscriber Token
   Details
```

Once incorporated into the blockchain the subscriber token is eligible for 
authoring new messages (not allowed in the same block nonetheless). Since by
construction any message must have a valid author, subscription is only on an
invitation basis. The first subscriber should be encoded in the gensis block
by the unique Eve message. The caption must be unique on the blockchain. Only
one message should be accepted for a caption string.

Details nonetheless can be modified by a request

```
About:
    Details
```

only to be accepted if the message Author token is the same as the subscription
token.

Audiences are created like

```
Create Audience:
   Token
   Description
```

If a valid subscriber is willing to be incorporated into an audiencea request 
to follow audience should be submitted in the form

```
Join Audience:
   Audience Token
   Expire Epoch
```

New members are added to the audience, or existing members are removed from
audience by

```
Audience Change:
   Audience Token
   Add: []Follower
   Remove: []Follower
   Details
```

where

```
Follower:
   Valid Subscriber Token
   Encryption by the Subscrber Token of the Private Key of the Audience Token
```
The message should only be accepted if every new follower can be traced back to 
a Join Audience request. Changes in audience can only be performed by the same 
author that created the audience.

The protocol ensures that a follower is capable of seeing the Tokens of all the 
other followers in the audience.

The follower can request to be removed from an audience by 

```
Withdraw Audience:
   Audience Token
```

This is only to be processed if the refered audience includes the message
author. 


An advertising offer is a message of the type

```
Advertising Offer:
   Offer Token
   Audience Token
   Content-Type
   Content-Data
   Advertising Fee
   Expire Epoch
   Signature by Advertising Wallet Token (Private Key) of the above
```

Expire Epoch must be greater than the epoch associated to the message. 
Advertising Fee will be withdrawn from fee message fee wallet upon valid
acceptance. 

The message to deliver content is of the form

```
Content:
   Audience Token
   Null or Encryption by audience token (private key) of a content key
   Content-Type
   Content-Data or Encryption using content key of Content-Data
   Null or Advertising Offer Token
```

The content can be public or encrypted. 

If the content reclaims an advertising offer the message is only valid if
there is an advertising offer with expire epoch before current epoch and content
type and content data are equal to the same information provided by the 
advertising offer. 

Finaly power of attorney is granted by

```
Grant power of attorney:
   Power of attorney token
   Expire Epoch
```

Power of attorney will allow the respective token to sign messages on behalf of
the author of the message granting the power of attorney, valid until a 
expire epoch. The power only becomes effective after the message is incorporated
into the blockchain.

In order to revoke the power of attorney before the expire epoch a message
must be submitted in the form

```
Revoke power of attorney:
   Power of attorney token
```

The message must be submitted by the author and signed by the author. 


Serialization
=============

Specific details about serialization of messages will be provided.

Cryptography
============

Specific details about cryptographic keys requirements will be provided.

Stability and Performance Estimations
=====================================

Once all the details of the protocol are nailed down estimations for the 
hardware, networkinging and storage resources to keep the network functioning
will be provided together with an assesment of its stability, its solidity and
its performance. 




