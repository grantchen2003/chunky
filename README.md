# chunky

## Reliable, resumable uploads for big files — because losing progress at 99% fucking sucks.

We've all been there. You're waiting for a massive file to upload, staring at the progress bar inching forward like a 3AM microwave. And just when it's about to finish, it fails. Maybe your internet cut out. Maybe your laptop died. Or maybe you fat-fingered `Ctrl-C` and terminated the upload. Whatever it was, you’re back to square one.

That’s where `chunky` comes in.

**`chunky` is a Go library for resumable uploads of arbitrarily large files. Now, if a large file upload fails halfway through, it picks up where it left off so you don’t have to start over from the beginning.**
