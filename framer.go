// 
// Copyright (c) 2014 Brian William Wolter, All rights reserved.
// Go Framer
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian William Wolter, Wolter Group New York, nor the
//     names of its contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package framer

import (
  "io"
  "fmt"
  "bytes"
  "encoding/binary"
)

const (
  SIZEOF_INT  = 4 // 4 byte int is used for frame length header (uint32)
  BUFFER_SIZE = 1024
)

/**
 * An input framer is responsible for segmenting distinct message frames as they are
 * read from an underlying input
 */
type InputFramer interface {
  
  /**
   * Read from the underlying input, possibly blocking until at least one message
   * is available. Zero or more frames are returned on success.
   */
  Read() ([][]byte, error)
  
}

/**
 * An output framer is responsible for framing distinct messages and writing them
 * to an underlying output.
 */
type OutputFramer interface {
  
  /**
   * Write a message to the underlying output, possibly blocking until the entire
   * message has been written.
   */
  Write([]byte) (error)
  
}

/**
 * An input framer that uses a reader as its underlying input
 */
type ReaderFramer struct {
  reader  io.Reader
  buffer  bytes.Buffer
}

/**
 * Create a reader framer
 */
func NewReaderFramer(reader io.Reader) *ReaderFramer {
  return &ReaderFramer{reader:reader}
}

/**
 * Read messages. This method blocks until at least one full message is read.
 */
func (r *ReaderFramer) Read() ([][]byte, error) {
  
  clen  := BUFFER_SIZE
  chunk := make([]byte, clen)
  
  for {
    
    // read from our input
    n, err := r.reader.Read(chunk)
    
    // append to our buffer
    if n > 0 {
      r.buffer.Write(chunk[:n]) // bytes.Buffer.Write() error is always nil according to the docs, so we don't check it
    }
    
    // check for a frame header
    if r.buffer.Len() > SIZEOF_INT {
      // if we have enough data to read at least one message, do so
      if flen := binary.BigEndian.Uint32(r.buffer.Bytes()); (r.buffer.Len() - SIZEOF_INT) >= int(flen) {
        return r.decode()
      }
    }
    
    // make sure nothing went wrong while reading (a reader may read valid bytes before producing an error)
    if err == io.EOF {
      return nil, io.EOF
    }else if err != nil {
      return nil, fmt.Errorf("Error reading from underlying input: %v", err)
    }
    
  }
  
  return nil, fmt.Errorf("A full frame could not be read")
}

/**
 * Decode frames from this framer's internal buffer. This method expects the buffer to
 * contain at least a single frame header (4 bytes), although not necessarily a full
 * message.
 */
func (r *ReaderFramer) decode() ([][]byte, error) {
  messages := make([][]byte, 0)
  
  for {
    
    // make sure we have at least one frame header
    if r.buffer.Len() < SIZEOF_INT {
      break // not enough data available for a header
    }
    
    // check the payload length
    flen := int(binary.BigEndian.Uint32(r.buffer.Bytes()))
    // make sure we have a full message in the buffer
    if r.buffer.Len() < SIZEOF_INT + flen {
      break // not enough data available for the entire frame
    }
    
    // consume the header
    r.buffer.Next(SIZEOF_INT)
    // set up our message buffer
    message := make([]byte, flen)
    // copy over bytes
    copy(message, r.buffer.Next(flen))
    // append our frame to the output set
    messages = append(messages, message)
    
    // loop to process more messages if we can...
  }
  
  return messages, nil
}

/**
 * An output framer that uses a Writer as its underlying input
 */
type WriterFramer struct {
  writer  io.Writer
}

/**
 * Create a writer framer
 */
func NewWriterFramer(writer io.Writer) *WriterFramer {
  return &WriterFramer{writer}
}

/**
 * Write a message to the underlying writer. This method blocks until the entire
 * message is written.
 */
func (w *WriterFramer) Write(message []byte) (error) {
  mlen := len(message)
  
  // write our header
  if err := binary.Write(w.writer, binary.BigEndian, uint32(mlen)); err != nil {
    return fmt.Errorf("Could not write message header: %v", err)
  }
  
  // write our message data
  for n := 0; n < mlen; {
    if z, err := w.writer.Write(message[n:]); err != nil {
      return fmt.Errorf("Could not write message data: %v", err)
    }else{
      n += z
    }
  }
  
  return nil
}

