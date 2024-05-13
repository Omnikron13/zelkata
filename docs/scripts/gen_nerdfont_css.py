import os
import re
import sys
import tempfile
import urllib.request
from urllib.error import HTTPError
import tarfile


# NerdFont release URL
BASE_URL = 'https://github.com/ryanoasis/nerd-fonts/releases/download/'

# There is sadly no 'latest' mechanism, although I could resort to strange git shenanigans to hoover up the font files
# directly from the master branch by, I believe, strange sparse checkout commands?
# Regardless a reasonably up-to-date version can be maintained here, and it can always be specified by argument too.
VERSION = 'v3.2.1'

# Currently the font archives are provided as either .zip or .tar.xz; the .tar.xz is notably smaller so the best choice
# right now, but available options may change in future.
ARCHIVE_TYPE = '.tar.xz'

# Best practice to get ourself a safe temporary directory to work in...
tmpDir = tempfile.mkdtemp()

# TODO: use 'argparse' to parse command line arguments for better... everything?
fontBaseName = sys.argv[1]
archiveName = f'{fontBaseName}{ARCHIVE_TYPE}'

# Override default version if argument passed
version = sys.argv[2] if len(sys.argv) > 2 else VERSION

# Full download URL for the specified font/version
url = f'{BASE_URL}{version}/{archiveName}'

# (attempt to) download the font archive
try:
   archivePath, _ = urllib.request.urlretrieve(url, os.path.join(tmpDir, archiveName))
except HTTPError:
   print(f'HTTP Error while attempting to download: {url} (bad version or font name?)')
   sys.exit(1)

# Suffixed name of the patched font
fontName = f'{fontBaseName}NerdFont'

# Output dir for our font files and generated CSS
fontDir = os.path.join(tmpDir, fontName)

# Pull just the actual font files from the archive
with tarfile.open(archivePath, 'r:xz') as tar:
   tar.extractall(fontDir, members=[m for m in tar.getmembers() if m.isfile() and m.name.endswith('.ttf')])

# We'll just name the CSS file the same as the font name, and store it in the font dir with the font files
cssFile = f'{fontName}.css'
cssPath = os.path.join(fontDir, cssFile)

# Generate @font-face entries for each file/variant provided in the given font
with open(os.path.join(fontDir, cssFile), 'w') as cssFile:
   cssFile.write(f'/* Generated from {fontName}-{version} by gen_nerdfont_css.py */\n\n')
   for f in os.listdir(fontDir):
      if f.endswith('.ttf'):
         weightRegex = r'(?P<weight>Thin|(?:Extra|Ultra)?Light|Normal|Regular|Medium|(?:Semi|Demi|Extra|Ultra)?Bold|Black|Heavy)'
         match = re.match(fr"^{fontName}(?P<spacing>Mono|Propo)?-{weightRegex}?(?P<stretch>Condensed)?(?P<style>Italic)?\.ttf$", f)

         spacing = ''
         if match.group('spacing') is not None:
            spacing = f' {match.group('spacing')}'

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


# Finally, output the path to the font directory
print(fontDir)

