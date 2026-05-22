# prik - make Braille characters on the fly

Braille characters are very useful, beyond, well, Braille. You see, in the
unicode spec, there is a block of Braille characters that are simple dot
patterns of 3x2 or 4x2. These do not have specific meaning attached to them, and
are very useful for all sorts of visual purposes.

This simple program shows a cool interactive screen where a user can create a
specific Braille character via `stdout`. Use the number keys to toggle the
(slightly wierdly) numbered dots, or use the arrow keys.

Because of this tool using `stdout` it is very easy to embed it or use it in
whatever weird workflows you have (I do not judge and I do not care).

## Example:

```bash
echo $(prik) >> somefile.txt
```
