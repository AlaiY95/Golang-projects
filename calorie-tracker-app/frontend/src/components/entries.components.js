import React, { useState, useEffect } from 'react'; // Import necessary dependencies from React
import axios from "axios"; // Import Axios for making HTTP requests
import { Button, Form, Container, Modal } from 'react-bootstrap'; // Import Bootstrap components
import Entry from './single-entry.component'; // Import a custom component

const Entries = () => {
    // Define and initialize various state variables using the useState hook
    const [entries, setEntries] = useState([]); // Array of calorie entries
    const [refreshData, setRefreshData] = useState(false); // Flag to trigger data refresh
    const [changeEntry, setChangeEntry] = useState({"change": false, "id": 0}); // State for changing an entry
    const [changeIngredient, setChangeIngredient] = useState({"change": false, "id": 0}); // State for changing ingredients
    const [newIngredientName, setNewIngredientName] = useState(""); // New ingredient name
    const [addNewEntry, setAddNewEntry] = useState(false); // Flag for adding a new calorie entry
    const [newEntry, setNewEntry] = useState({"dish": "", "ingredients": "", "calories": 0, fat: 0}); // New calorie entry data


     // useEffect hook to fetch all entries when the component mounts
    useEffect(() => {
        getAllEntries();
    }, [])

    // Refresh data when the refreshData flag is true
    if(refreshData){
        setRefreshData(false);
        getAllEntries();
    }

    return(
        <div>
            <Container>
        <Button onClick={() => setAddNewEntry(true)}>Track today's calories</Button>
            </Container>
            <Container>
        {/* Map over the entries and render Entry components for each */}
        {entries != null && entries.map((entry, i) =>(
            <Entry entryData={entry} deleteSingleEntry={deleteSingleEntry} setChangeIngredient={setChangeIngredient} setChangeEntry={setChangeEntry} />
        ))}
            </Container>

            {/* Modal for adding a new calorie entry */}
            <Modal show={addNewEntry} onHide={() => setAddNewEntry(false)} centred>
            <Modal.Header closeButton>
            <Modal.Title>Add Calorie Entry</Modal.Title>
            </Modal.Header>

            <Modal.Body>
                {/* Form for entering new calorie entry details */}
                <Form.Group>
                    <Form.Label>dish</Form.Label>
                    <Form.Control onChange={(event) => {newEntry.dish = event.target.value}}></Form.Control>
                    <Form.Label>ingredients</Form.Label>
                    <Form.Control onChange={(event) => {newEntry.ingredients = event.target.value}}></Form.Control>
                    <Form.Label>calories</Form.Label>
                    <Form.Control onChange={(event) => {newEntry.calories = event.target.value}}></Form.Control>
                    <Form.Label>fat</Form.Label>
                    <Form.Control type="number" onChange={(event) => {newEntry.fat = event.target.value}}></Form.Control>
                </Form.Group>
                <Button onClick={() => addSingleEntry()}>Add</Button>
                <Button onClick={()=> setAddNewEntry(false)}>Cancel</Button>
            </Modal.Body>
            </Modal>

            {/* Modal for changing ingredients */}
            <Modal show={changeIngredient.change} onHide={() => setChangeIngredient({"change": false, "id":0})} centred>
            <Modal.Header closeButton>
                <Modal.Title>Change Ingredients</Modal.Title>
            </Modal.Header>

            <Modal.Body>
                <Form.Group>
                    <Form.Label>new ingredients</Form.Label>
                    <Form.Control onChange={(event) => {setNewIngredientName(event.target.value)}}></Form.Control>
                </Form.Group>
                <Button onClick={() => changeIngredientForEntry()}>Change</Button>
                <Button onClick={() => setChangeIngredient({"change": false, "id":0})}>Cancel</Button>
            </Modal.Body>
            </Modal>

            {/* Modal for changing an entry */}
            <Modal show={changeEntry.change} onHide={() => setChangeEntry({"change": false, "id":0})} centred>
            <Modal.Header closeButton>
                <Modal.Title>Change Entry</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                {/* Form for changing entry details */}
                <Form.Group>
                    <Form.Label>dish</Form.Label>
                    <Form.Control onChange={(event) => {newEntry.dish = event.target.value}}></Form.Control>
                    <Form.Label>ingredients</Form.Label>
                    <Form.Control onChange={(event) => {newEntry.ingredients = event.target.value}}></Form.Control>
                    <Form.Label>calorie</Form.Label>
                    <Form.Control onChange={(event) => {newEntry.calories = event.target.value}}></Form.Control>
                    <Form.Label>fat</Form.Label>
                    <Form.Control type="number" onChange={(event) => {newEntry.fat = event.target.value}}></Form.Control>
                </Form.Group>
                <Button onClick={() => changeSingleEntry()}>Change</Button>
                <Button onClick={() => setChangeEntry({"change": false, "id":0})}>Cancel</Button>
            </Modal.Body>
        </Modal>
        </div> 
    );

    // Responsible for making an HTTP PUT request to update the ingredients for a specific entry. 
    function changeIngredientForEntry(){
        //  It sets the changeIngredient flag to false, which typically controls the visibility of a modal or form
        changeIngredient.change = false
        // Constructs the URL for the update request, including the changeIngredient.id
        var url = "http://localhost:8000/ingredient/update/" + changeIngredient.id
        // Sends an axios PUT request to the server with the updated ingredient data
        axios.put(url, {
            "ingredients": newIngredientName
        }).then(response => {
            console.log(response.status)
            // Upon a successful response (HTTP status 200), it sets the refreshData flag to true
            if(response.status == 200 ){
                setRefreshData(true)
            }
        })
    }

    function changeSingleEntry(){
        changeEntry.change = false;
        var url = "http://localhost:8000/entry/update/" + changeEntry.id
        axios.put(url, newEntry)
        .then(response =>{
            if(response.status == 200){
                setRefreshData(true)
            }
        })
    }

    function addSingleEntry(){
        setAddNewEntry(false)
        var url = "http://localhost:8000/entry/create"
        axios.post(url, {
            "ingredients":newEntry.ingredients,
            "dish": newEntry.dish,
            "calories": newEntry.calories,
            "fat": parseFloat(newEntry.fat)
        }).then(response => {
            if(response.status == 200){
                setRefreshData(true)
            }
        })
    }
    
    function deleteSingleEntry(id){
        var url = "http://localhost:8000/entry/delete/" + id
        axios.delete(url, {
    
        }).then(response => {
            if (response.status == 200){
                setRefreshData(true)
            }
        })
    }
    
    function getAllEntries(){
        var url = "http://localhost:8000/entries"
        axios.get(url, {
            reponseType: 'json'
        }).then(response => {
            if(response.status == 200){
                setEntries(response.data)
            }
        })
    }
}

export default Entries
