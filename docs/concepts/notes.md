󰭷 Notes
=======

The main building blocks of _Zelkata_ are _󰭷 notes._
A _󰭷 note_ is, in general, some information about a single concept.
It's a short piece of information, a thought, a reminder, a project idea, a plan, etc.
In _Zelkata_, as with the _Zettelkästen_ systems that it is based upon, emphasis is placed on the importance of keeping
notes short and specific; you can always add more notes, and if you find a note is becoming unfocussed or covering too
much, you should really split it into multiple notes.

NOTE: The tools to perform actions like splitting notes are not currently available, but are planned for the near
future. Again, a CLI option will likely come first, followed up ASAP by more integrated and user-friendly TUI options.

The actual on-disk format of a _󰭷 note_, not that you should generally have to access raw _󰭷 note_ files directly, is:

```yaml
---
[YAML frontmatter/meta-data]
...

```commonmark
[Markdown content]
```

