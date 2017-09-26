# fix-mp3-tag
A lot of Cyrillic mp3's are broken by f%&amp;king windows encoding. Trying to fix this. 

You need id3info and id3tag applications, usually they are part of id3lib package for most Linux distribution.

To install:
* setup the environment:

  ```
  mkdir -p devel/golang/{bin,src}
  export GOPATH=$PWD/devel/golang
  ```

* bring the package from github:

  ```
  cd $GOPATH
  git clone https://github.com/bukind/fix-mp3-tag.git src/bukind/fix-mp3-tag
  ```

* install the binary:

  ```
  go install github.com/bukind/fix-mp3-tag
  ```
  
  The binary will be created as $GOPATH/bin/fix-mp3-tag.
