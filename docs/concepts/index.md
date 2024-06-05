󱥬 Concepts
==========

There are only a few core concepts in _Zelkata_, all of which it should be easy to grasp the basics of and begin
working with.
This isn't to say that there isn't depth to each and every once of these concepts and systems, but the design intent is
for their depth to be either emergent or optional, as a core design principle in _Zelkata_ is to minimise friction as
much as possible for you.

NOTE: The above may not be _entirely_ applicable to the early tools and methods of interacting with _Zelkata_,
as they more user-friendly and intuitive tools are going to progressively follow on from , to borrow a term from the
development of Git itself, the _'plumbing'_ of the system. The progression is likely to be be low level CLI tools ->
much more pleasant TUI interfaces for those happy in modern terminal environments -> GUI/Web/App interfaces for those
not as keen living in the terminal.


󰭷 Notes
-------

The main building blocks of _Zelkata_ are _󰭷 notes._
A _󰭷 note_ is, in general, some information about a single concept.
It's a short piece of information, a thought, a reminder, a project idea, a plan, etc.
In _Zelkata_, as with the _Zettelkästen_ systems that it is based upon, emphasis is placed on the importance of keeping
notes short and specific; you can always add more notes, and if you find a note is becoming unfocussed or covering too
much, you should really split it into multiple notes.

NOTE: The tools to perform actions like splitting notes are not currently available, but are planned for the near
future. Again, a CLI option will likely come first, followed up ASAP by more integrated and user-friendly TUI options.

Notably more in-depth information about [_󰭷 Notes_](concepts/notes.md) can be found on their own page, both some
slightly more technical information about how they work, and tips, best practices, examples of use, and more advanced
optional features.


 Tags
------

If _󰭷 notes_ are the bricks of _Zelkata_, then [_ Tags_](concepts/tags.md) could be considered the mortar.
At the most basic, a _ tag_ is a what allows you to split your _󰭷 notes_ on a particular subject into those smaller
and more focussed 

A _󰭷 note_ briefly describing a 🪈flute would likely be given the _ tag_ 'Instrument', as would a _ note_ describing a
🪕banjo
These _󰭷 notes_ seem obviously separate but related, whereas a _󰭷 note_ describing an _󰋄 electric guitar_ might get
entangled with information that ought to be kept in an _󰝱 acoustic guitar_ _󰭷 note_, which is where you would split
them and rely on the 'Instrument' _ tag_ to easily find these related items by browsing your 'Instrument' _ tags_.

Because _ tags_ are very much not folders/directories, and you may have as many arbitrary _ tags_ as you like, you
would likely add, say, a 'Stringed' _ tag_ to both of them, even a 'Guitar' _ tag_, allowing you to cross-reference
or filter your _󰭷 notes_ in a variety of ways.

It should seem obvious at this point that you would probably then add _󰓹 tags_ like 'Music' to them, maybe 'Hobbies',
etc.

Although you only have to enter arbitrary _ tag_ names when writing notes, it should be pretty frictionless to
classify your notes at the time, but the number of _ tags_ the system will be starting to keep automatic track of for
you will start to grow pretty fast and become a little nebulous.
You will want to perform some manual curation of your _ tags_ as this starts to happen, which isn't difficult but
begins to impose more structure to your _󱟱 notes collection_, and open up more options for exploring, filtering, and
cross-referencing your notes.
The system will even begin to be able to start bringing up notes that you may not have even really considered the
relations between.

Details on [_ Tags_](concepts/tags.md) can be found on their own page, which will go into more depth on how you can
use them in a more _Zelkata_/_Zettelkästen_, less 'flat' way.
You should learn some of the technical details of how _ tags_ work, and more importantly how you can get the most out
of them.

NOTE: You will also learn about concepts such as _󰓼 virtual tags_, which don't actually behave in the basic way this
section describes, and potentially go against some of the principles which this page has tried to keep simple.

