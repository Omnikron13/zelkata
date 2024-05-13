import os
import sys
import re


# TODO: use 'argparse' to parse command line arguments for better... everything?
fontDir = sys.argv[1]
fontBaseName = os.path.basename(fontDir)
fontName = fontBaseName + "NerdFont"

cssFile = f'{fontName}.css'

with open(os.path.join(fontDir, cssFile), 'w') as cssFile:
   for f in os.listdir(fontDir):
      if f.endswith('.ttf'):
         weightRegex = r'(?P<weight>Thin|(?:Extra|Ultra)?Light|Normal|Regular|Medium|(?:Semi|Demi|Extra|Ultra)?Bold|Black|Heavy)'
         match = re.match(fr"^{fontName}(?P<spacing>Mono|Propo)?-{weightRegex}?(?P<stretch>Condensed)?(?P<style>Italic)?\.ttf$", f)

         spacing = ''
         if match.group('spacing') is not None:
             spacing = ' ' + match.group('spacing')

         weight = 400
         if match.group('weight') is not None:
            weight = {
               'Thin': 100,
               'ExtraLight': 200,
               'UltraLight': 200,
               'Light': 300,
               'Normal': 400,
               'Regular': 400,
               'Medium': 500,
               'SemiBold': 600,
               'DemiBold': 600,
               'Bold': 700,
               'ExtraBold': 800,
               'UltraBold': 800,
               'Black': 900,
               'Heavy': 900,
            }[match.group('weight')]

         stretch = 'normal'
         if match.group('stretch') is not None:
             stretch = match.group('stretch')

         style = 'normal'
         if match.group('style') is not None:
             style = match.group('style')

         cssFile.write(f'@font-face {{\n')
         cssFile.write(f'  font-family: "{fontName}{spacing}";\n')
         cssFile.write(f'  src:\n   url("{f}") format("truetype");\n')
         cssFile.write(f'  font-weight: {str(weight)};\n')
         cssFile.write(f'  font-style: {style};\n')
         cssFile.write(f'  font-stretch: {stretch};\n')
         cssFile.write(f'}}\n\n')

   cssFile.write(f'\n:root {{\n')
   cssFile.write(f'  --md-text-font: "{fontName} Propo";\n')
   cssFile.write(f'  --md-code-font: "{fontName}" "{fontName} Mono";\n')
   cssFile.write(f'}}\n\n')

