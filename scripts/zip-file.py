#!/usr/bin/env python
import os
from os.path import isdir
import argparse
import zipfile
import hashlib


def zip_asset(file,destination,arcname,version,goos,goarch):
  if not isdir(destination):
    os.mkdir(destination)

  filename = os.path.basename(file)
  output = "{0}/{1}_{2}_{3}_{4}.zip".format(destination,arcname,version,goos,goarch)

  with zipfile.ZipFile(output,'w') as f:
    f.write(filename=file,arcname=arcname)
    f.comment=filename
  return output


if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument("--prefix", nargs='+', help="prefix strings for targets")
  parser.add_argument("--destination", default="release", help="Destination folder for files")
  parser.add_argument("--version", default=os.environ.get('VERSION'), help="Version number for binary, e.g.: v1.0.0")
  parser.add_argument("--goos", default=os.environ.get('GOOS'), help="GOOS parameter")
  parser.add_argument("--goarch", default=os.environ.get('GOARCH'), help="GOARCH parameter")
  args = parser.parse_args()

  if args.version is None:
    raise parser.error("argument --version is required")
  if args.goos is None:
    raise parser.error("argument --goos is required")
  if args.goarch is None:
    raise parser.error("argument --goarch is required")

  for prefix in args.prefix:
    fname = "build/{0}_{1}_{2}".format(prefix, args.goos,args.goarch)
    file = zip_asset(fname,args.destination,prefix,args.version,args.goos,args.goarch)
    print(file)
