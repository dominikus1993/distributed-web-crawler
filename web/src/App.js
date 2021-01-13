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
  const [url, setUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
 
  useEffect(() => {
    const fetchData = async () => {
      setIsError(false);
      setIsLoading(true);
 
      try {
        const result = await axios(url);
 
        setData(result.data);
      } catch (error) {
        setIsError(true);
      }
 
      setIsLoading(false);
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
          setUrl(`http://hn.algolia.com/api/v1/search?query=${query}`)
        }
      >
        Search
      </button>
  </form>
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