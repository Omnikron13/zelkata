// Description: This script adds the classes 'accent' and 'nerd-font' to the first character of the selected elements,
//              if said character is part of the NerdFont glyph set.

window.addEventListener('DOMContentLoaded', () => {
   const accentify = [
     ...document.querySelectorAll('h1, h2, a, em, .md-ellipsis'),
   ]

   const nerdRanges = [
      {name: 'Octicons', ranges: [{first: 0xf400, last: 0xf533}, {first: 0x2665, last: 0x2665}, {first: 0x26A1, last: 0x26A1}]},
      {name: 'IEC Power Symbols', ranges: [{first: 0x23fb, last: 0x23fe}, {first: 0x2b58, last: 0x2b58}]},
      {name: 'Pomicons', ranges: [{first: 0xe000, last: 0xe00a}]},
      {name: 'Powerline Symbols', ranges: [{first: 0xe0a0, last: 0xe0a2}, {first: 0xe0b0, last: 0xe0b3}]},
      {name: 'Powerline Extra Symbols', ranges: [{first: 0xe0a3, last: 0xe0a3}, {first: 0xe0b4, last: 0xe0c8}, {first: 0xe0ca, last: 0xe0ca}, {first: 0xe0cc, last: 0xe0d7}]},
      {name: 'Font Awesome Extension', ranges: [{first: 0xe200, last: 0xe2a9}]},
      {name: 'Weather Icons', ranges: [{first: 0xe300, last: 0xe3e3}]},
      {name: 'Seti-UI + Custom', ranges: [{first: 0xe5fa, last: 0xe6b5}]},
      {name: 'Devicons', ranges: [{first: 0xe700, last: 0xe7c5}]},
      {name: 'Codicons', ranges: [{first: 0xea60, last: 0xec1e}]},
      {name: 'Font Awesome', ranges: [{first: 0xed00, last: 0xefce}, {first: 0xf000, last: 0xf2ff}]},
      {name: 'Font Logos/Font Linux', ranges: [{first: 0xf300, last: 0xf372}]},
      {name: 'Material Design Icons', ranges: [{first: 0xf0001, last: 0xf1af0}]},
   ]

   accentify.forEach((el) => {
      const text = el.innerText
      const firstChar = text.codePointAt(0)
      if (nerdRanges.map(r => r.ranges).flat().some(r => firstChar >= r.first && firstChar <= r.last)) {
         const glyph = String.fromCodePoint(firstChar)
         const tail = text.substring(glyph.length)
         const accentSpan = document.createElement('span')
         accentSpan.classList.add('accent', 'nerd-font')
         accentSpan.innerText = glyph
         el.replaceChildren(accentSpan, tail)
      }
   })
})
