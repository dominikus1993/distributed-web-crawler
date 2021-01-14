import React, { useState, useEffect } from "react";
import socketIOClient from "socket.io-client";
import axios from "axios"
const ENDPOINT = "http://127.0.0.1:4001";


/**
 * @param {{ url: string, content: {url: string}[] }} props
 */
function ContentList(props) {
  return (
    <div>
      <p>Url: {props.url}</p>
      <ul>
        {props.contents.map(content => (
          <li key={content.url}>
            <img src={content.url} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function ParseUrlForm() {
  const [url, setUrl] = useState("");
  const [query, setQuery] = useState("")

  useEffect(() => {
    const fetchData = async () => {
      if (url !== "") {
        console.log("Request");
        await axios.post(ENDPOINT, { url: url });
      }

    };
    fetchData();
  }, [url]);
  return <form>
    <input
      type="text"
      value={query}
      onChange={event => setQuery(event.target.value)}
    />
    <button
      type="button"
      onClick={() =>
        setUrl(query)
      }
    >
      Search
      </button>
  </form>
}

function App() {
  const [response, setResponse] = useState({ hasData: false, url: "", contents: [] });

  useEffect(() => {
    const socket = socketIOClient(ENDPOINT);
    socket.on("new-crawled-media", data => {
      setResponse({ hasData: true, url: data.url, contents: data.contents });
    });
    return () => socket.disconnect()
  }, []);

  if (response.hasData) {
    return <ContentList url={response.url} contents={response.contents} ></ContentList>
  }
  else {
    return <ParseUrlForm></ParseUrlForm>
  }
}

export default App;