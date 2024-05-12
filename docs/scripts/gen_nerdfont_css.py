import os
import sys
import re


# TODO: use 'argparse' to parse command line arguments for better... everything
print(sys.argv)

fontDir = sys.argv[1]
fontBaseName = os.path.basename(fontDir)
fontName = fontBaseName + "NerdFont"

ccsDir = sys.argv[2]
cssFile = fontName + '.css'
cssPath = os.path.join(ccsDir, cssFile)

fontAbsDir = sys.argv[3]

with open(cssPath, 'w') as cssFile:
   for f in os.listdir(fontDir):
      if f.endswith('.ttf'):
         match = re.match(r'^.+?(?P<spacing>Mono|Propo)?-(?P<weight>Thin|(?:Extra|Ultra)?Light|Normal|Regular|Medium|(?:Semi|Demi|Extra|Ultra)?Bold|Black|Heavy)?(?P<stretch>Condensed)?(?P<style>Italic)?\.ttf$', f)
         match = re.match(r'^.+?(?P<spacing>Mono|Propo)?-(?P<weight>Thin|(?:Extra|Ultra)?Light|Normal|Regular|Medium|(?:Semi|Demi|Extra|Ultra)?Bold|Black|Heavy)?(?P<stretch>Condensed)?(?P<style>Italic)?.+?$', f)
         #print(match)

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

         print(spacing, weight, stretch, style)

         #fontName = f.split('.')[0]
         cssFile.write('\n@font-face {\n')
         cssFile.write('  font-family: "' + fontName + spacing + '";\n')
         # TODO: reowr this to be a lot less fragile
         cssFile.write('  src:\n   url("'  + fontAbsDir + fontBaseName + '/' + f + '") format("truetype");\n')
         cssFile.write('  font-weight: ' + str(weight) + ';\n')
         cssFile.write('  font-style: ' + style + ';\n')
         cssFile.write('  font-stretch: ' + stretch + ';\n')
         cssFile.write('}\n')
         cssFile.write('\n')

   cssFile.write('\n:root {\n')
   cssFile.write('  --md-text-font: "' + fontName + ' Propo";\n')
   cssFile.write('  --md-code-font: "' + fontName + '";\n')
   cssFile.write('}\n')

