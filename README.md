socklog
=======

socklog forward TCP port to another ADDR:PORT and dump stream.

## Usage

* Listen 0.0.0.0:8100. When incoming connection established, socklog try to connect 10.0.0.1:80 and forward stream and dump it to stdout.  
  Control code(means 0x00-0x1f except 0x09, 0x0a, 0x0c, 0x0d) is masked using '.'
  ```
  socklog --bind 0.0.0.0:8100 --dest 10.0.0.1:80
  ```
* Control code is not masked
  ```
  socklog --bind 0.0.0.0:8100 --dest 10.0.0.1:80 --mask=false
  ```
* Dump stream to the file instead of stdout.
  ```
  socklog --bind 0.0.0.0:8100 --dest 10.0.0.1:80 --out /path/to/file
  ```

## My environment

* macOS Sierra
  * go 1.7.1
 