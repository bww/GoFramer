// 
// Copyright (c) 2014 Brian W. Wolter, All rights reserved.
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
// Developed in New York City
// 

package framer

import (
  "io"
  "bytes"
  "encoding/binary"
)

const (
  SIZEOF_INT  = 4 // 4 byte int is used for frame length header (uint32)
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
 * Read messages. This method blocks until at least one full message is read.
 */
func (r *ReaderFramer) Read() ([][]byte, error) {
  
  clen  := 24 // TEMPORARY! make this larger
  chunk := make([]byte, blen)
  
  for {
    
    // read from our input
    n, err := r.reader.Read(chunk)
    
    // append to our buffer
    if n > 0 {
      r.buffer.Write(chunk[:n]) // bytes.Buffer.Write() error is always nil according to the docs, so we don't check it
    }
    
    // make sure nothing went wrong while reading (a reader may read valid bytes before producing an error)
    if err != nil {
      return nil, fmt.Errorf("Error reading from underlying input: %v", err)
    }
    
    // check for a frame header
    if r.buffer.Len() > SIZEOF_INT {
      // if we have enough data to read at least one message, do so
      if mlen := binary.BigEndian.Uint32(r.buffer.Bytes()); r.buffer.Len() >= mlen {
        return r.decode()
      }
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
  messages := make([][]byte)
  
  for {
    var flen uint32
    
    // make sure we have at least one frame header
    if r.buffer.Len() < SIZEOF_INT {
      break
    }
    
    // check the header length
    if flen = binary.BigEndian.Uint32(r.buffer.Bytes()); r.buffer.Len() < flen {
      return nil, fmt.Println("Could not decode frame header: %v", err)
    }
    
    // skip the header
    r.buffer.Next(SIZEOF_INT)
    
    // set up our message buffer
    message := make([]byte, flen)
    
    // read our message data
    if n, err := r.buffer.Read(message); n < flen {
      return nil, fmt.Println("Could not read entire frame: %d < %d", n, flen)
    }else if err != nil {
      return nil, fmt.Println("Could not read entire frame: %v", err)
    }
    
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
  
}


