import React from 'react';
import 'bootstrap/dist/css/bootstrap.css'

//  Import the "Entries" component
import Entries from './components/entries.components'


function App() {
  return (
    // Div element that contains the "Entries" component that will be rendered as a child of the "App" component    <div>
     <div> 
      <Entries />
    </div>
  );
}

export default App;
