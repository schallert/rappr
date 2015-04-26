# Rappr

This is a small silly app for tallying votes in the epic Derek vs. Jay rap battle of 2015. It simply maintains a tally of two votes (whether the text message voted for "Jay" or "Derek") that come in via Twilio SMS.

**THIS IS JUST A SMALL APP THAT WAS THROWN TOGETHER QUICKLY** as a chance to try out [BoltDB](https://github.com/boltdb/bolt) and throw something together with Go. It is not secure (requests can come from anywhere, there's no authentication that it comes from Twilio). There's no guarantee that error checking is robust. There's no guarantee that this won't blow up your entire server. You get the idea.
