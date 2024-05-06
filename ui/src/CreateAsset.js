import React, { useState } from 'react';
import { Button, TextField } from '@material-ui/core';
import axios from 'axios';

const CreateAsset = () => {
  const [name, setName] = useState('');
  const [price, setPrice] = useState(0);

  const handleCreate = async () => {
    try {
      const response = await axios.post('/api/assets', { name, price });
      console.log('Asset created:', response.data);
      // Reset form fields
      setName('');
      setPrice(0);
    } catch (error) {
      console.error('Error creating asset:', error);
    }
  };

  return (
    <div>
      <TextField
        label="Asset Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <TextField
        label="Price"
        type="number"
        value={price}
        onChange={(e) => setPrice(parseFloat(e.target.value))}
      />
      <Button variant="contained" color="primary" onClick={handleCreate}>
        Create Asset
      </Button>
    </div>
  );
};

export default CreateAsset;