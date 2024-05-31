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

The actual on-disk format of a _󰭷 note_, not that you should generally have to access raw _󰭷 note_ files directly,
begins with a YAML frontmatter/meta-data block, followed by the main contents of the note itself:

```yaml
---
id: [generually a UUID]
tags: [list of tags]
...

```markdown
Generally MarkDown
==================

The main contents of the note will generally , or at least by default, be in MarkDown format.

A particular flavour of MarkDown can be set in your configuration file, depending on your preferences.

You can also, if you wish, set any format you like essentially in the configuration file; AsciiDoc, ReStructuredText,
LaTeX, etc.
It is likely easiest to leave the defaults unless you are more familiar with another format and prefer it.
```

