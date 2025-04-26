# syncfo

> [!note]
> Currently this repo is for gathering ideas, brainstorming, and planning out the final app. There's no usable app as of now. If you have ideas, don't hesitate to create issues where we can discuss them.

## The Problem

More and more communities of open-source tools move away from web forums/ stackoverflow over to Discord servers.

Why is that bad? Discord was primarly made for chatting and not for information retention. Discussions in support servers constantly face the problem of the same questions turning up over and over, just because any previous discussions are hard to find and you also can't find anything when googling in your browser.

Set yourself in the shoes of somebody starting with some programming library you have no experience with (for example the Discord.py library). The documentation might be really good, but then you stumble upon some problem you can't find anything in the docs about. Or say you're trying to implement something that the library doesn't quite make clear how it's supposed to be done.

Nowadays, these are a fairly logical set of steps to take to answer your desires:

1. Ask your LLM of choice to explain how to implement your wish. If you're lucky, the library has been used quite a bit and so a lot of it is in the training set. Then the LLM is able to nudge you in the right direction and off you go. With more niche libraries you'll not get much helpful info unless the problem is more on the trivial side.
2. Google it! Search engines are really good at doing their job: indexing the web and trying to match up your search query with the most relevant places. Some libraries you'll be able to find answers in threads on the library's own forum. For more popular libraries you might even be lucky and find discussions on stackoverflow. And on less popular libraries you might even stumble upon some github issues. But with a large part of support discussions happening on Discord, which search engines can't index, you're now left with:
3. Head to the library's Discord server. Often there are some support channels you can scoure through to find what you need. Or you can start chatting and hope some kind soul is able to help you out.

I've personally even started jumping straight to 3. because I already know the library doesn't have much presence on the web.

Just relying on a Discord server sounds all nice and dandy. It's easy to set up a server after all. But there are a two major problems:

- What if the Discord server closes or more likely a support channel is hidden or even deleted? All the discussions and information is suddenly just gone. On a web forum even if the forum goes down you might be able to find threads in the waybackmachine or have snippets displayed by the search engine.
- The Discord search is fairly bad. Search engines have been optimized for exactly this; to be extremely good at searching. Take a look at stackoverflow; most people access stackoverflow threads straight from google and don't even attempt searching through the SO website directly.

## Comparing alternatives

### Classical web forums

```asciidoc
- hard to setup
+ indexed by search engines
```

### Discourse

> [!note]
> todo

## syncfo

The general idea with syncfo is to sync Discord support channels/forums to a (simple) web forum. This is to allow for discussions to still be kept on Discord, but by syncing everything to a web forum it is possible to get threads indexed by search engines.

## Additional ideas

Just a few other brainstormed ideas.

### No web application

The web forum should be a fairly static website and not a web application. Basically meaning every page is accessibly by simply visiting the relevant URL and there's no javascript required to load certain information. This is a key requirement to make a forum be indexable.

### Github pages compatible

It should be very easy to host the website and not necessairily require a dedicated server to host the forum. A lot of open-source communities use github pages to host their documentation for example.

That would mean there can't be a proper database or backend tools for hosting the web forum itself. Syncing itself would be done by a Discord bot or similar.
Something to look into would be [utterances](https://utteranc.es/): A tool that stores comments in github issues so a website doesn't require an explicit database. A potential problem could be that depending on github issues could become very slow, abusable, and make everything a lot more complicated than it has to be.

### Name anonymization

One question is how to depict messages on the web forum. Should names and profile pictures be carried over or should names be randomly generated to anonymize discord messages? In such forums the names themselves are rarely of importance. By using easy-to-remember aliases it would be possible to distinguish messages from the OP without putting up all usernames online.

### Store images separately

A common problem with old forums is that images are linked from some third-party source like imgur or elsewhere, causing a lot of images to simply not exist anymore. A lot of image links on old forums are dead. Is the solution for this to store any images locally and link those, similarily to how discord does it? What are the problems other that eating up storage faster?
