import argparse
import os
import re
import shutil
import sys
import tarfile
import tempfile
import urllib.request
from urllib.error import HTTPError


# NerdFont release URL
BASE_URL = 'https://github.com/ryanoasis/nerd-fonts/releases/download/'

# There is sadly no 'latest' mechanism, although I could resort to strange git shenanigans to hoover up the font files
# directly from the master branch by, I believe, strange sparse checkout commands?
# Regardless a reasonably up-to-date version can be maintained here, and it can always be specified by argument too.
VERSION = 'v3.2.1'

# Currently the font archives are provided as either .zip or .tar.xz; the .tar.xz is notably smaller so the best choice
# right now, but available options may change in future.
ARCHIVE_TYPE = '.tar.xz'

# I believe all the 'official' patched font files are TTF files
FONT_EXTENSION = 'ttf'

# Parse our command line arguments like a real live boy
argParser = argparse.ArgumentParser(
   prog="nerdfontweb.py",
   description="Fetch & process patched Nerd Fonts, generating required CSS to use them as web fonts",
   epilog="Originally part of Zelkata: https://github.com/omnikron13/zelkata",
)
argParser.add_argument('archive', metavar='name', help="name of the font (or at least the archive) to fetch, sans 'NerdFont'")
argParser.add_argument('-v', '--version', default=VERSION, help="fully qualified release version (including 'v' prefix) to fetch, if the default is out-of-date")
argParser.add_argument('-o', '--output', default=None, help='directory path to move output files/directories to (implies --clean)')
argParser.add_argument('-n', '--name', default=None, help='(base) name of font, if it differs from that of the archive file')
argParser.add_argument('-q', '--quiet', action='store_true', help='supress output of temporary output directory path (implies --clean)')
argParser.add_argument('-c', '--clean', action='store_true', help='delete temporary working dir when done (implies --quiet)')
argParser.add_argument('--gen-default', action='store_true', help='generate additional `default-font.css` file including the primary CSS file & setting the processed font as the default')
argParser.add_argument('--mkdocs', action='store_true', help='set default value for --output suitable for MkDocs, and set MkDocs specific variables in `default-font.css` if generated (implies --output)')
# TODO: control behaviour on failure to parse font file name
# TODO: alternate weight conversion table/system?

args = argParser.parse_args()
#print(args)

# Process the implications listed in the help text
if args.mkdocs and args.output is None:
   args.output = 'docs/assets/fonts/'
if args.quiet:
   args.clean = True
if args.clean:
   args.quiet = True
if args.output is not None:
   args.clean = True

# Best practice to get ourself a safe temporary directory to work in...
tmpDir = tempfile.mkdtemp()

# Full filename for the archive to fetch
archiveName = f'{args.archive}{ARCHIVE_TYPE}'

# Name of font sans 'NerdFont' suffix
fontBaseName = args.archive if args.name is None else args.name

# Override default version if argument passed
version = args.version if args.version is not None else VERSION

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
   tar.extractall(fontDir, members=[m for m in tar.getmembers() if m.isfile() and m.name.endswith(f'.{FONT_EXTENSION}')])

# We'll just name the CSS file the same as the font name, and store it in the font dir with the font files
cssFile = f'{fontName}.css'
cssPath = os.path.join(fontDir, cssFile)

# Generate @font-face entries for each file/variant provided in the given font
with open(os.path.join(fontDir, cssFile), 'w') as cssFile:
   cssFile.write(f'/* Generated from {fontName}-{version} by gen_nerdfont_css.py */\n\n')
   for f in os.listdir(fontDir):
      if f.endswith('.ttf'):
         weightRegex = r'(?P<weight>Thin|(?:Extra|Ultra|Semi)?Light|Normal|Regular|Medium|(?:Semi|Demi|Extra|Ultra)?Bold|(?:Extra|Ultra)?Black|Heavy)'
         match = re.match(fr"^(?P<name>.+NerdFont)(?P<spacing>Mono|Propo)?-{weightRegex}?(?P<stretch>Condensed)?(?P<style>Italic)?\.{FONT_EXTENSION}$", f)

         # Fail completely if a font filename can't be processed. Perhaps a bit excessive...
         if match is None:
            print(f"Failed to match/extract font file properties for: {f}")
            # TODO: control this behaviour with a cmd line flag
            sys.exit(1)

         spacing = ''
         if match.group('spacing') is not None:
            spacing = f' {match.group('spacing')}'

         weight = 400
         if match.group('weight') is not None:
            weight = {
               'Thin': 100,
               'Hairline': 100,
               'ExtraLight': 200,
               'UltraLight': 200,
               'Light': 300,
               'SemiLight': 350,
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
               'ExtraBlack': 950,
               'UltraBlack': 950,
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


if not args.quiet:
   # Finally, output the path to the font directory
   print(fontDir)

