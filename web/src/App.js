import React, { useState, useEffect } from "react";
import socketIOClient from "socket.io-client";
const ENDPOINT = "http://127.0.0.1:4001";

function list(data) {
  return (
    <div>
      <p>Url: {data.url}</p>
      <ul>
        {data.contents.map(content => (
          <li key={content.url}>
            <img src={content.url} />
          </li>
        ))}
      </ul>
    </div> 
  );
}

function form() {

}

function App() {
  const [response, setResponse] = useState({hasData: false, url: "", contents: [] });
  
  useEffect(() => {
    const socket = socketIOClient(ENDPOINT);
    socket.on("new-crawled-media", data => {
      setResponse({hasData: true, url: data.url, contents: data.contents });
    });
    return () => socket.disconnect()
  }, []);

  if(response.hasData) {
    return list(response)
  }
  else {
    return <h1>Witam</h1>
  }
}

export default App;