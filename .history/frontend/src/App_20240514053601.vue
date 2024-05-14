<template>
  <div class="container">
    <h1>Assets Portfolio Management</h1>
    <div class="chart-container">
      <canvas ref="doughnutChart"></canvas>
    </div>
    <table>
      <thead>
        <tr>
          <th>Asset Symbol</th>
          <th>Quantity</th>
          <th>Current Market Price</th>
          <th>Total Value</th>
          <th>Option</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(asset, index) in assets" :key="index">
          <td>{{ asset.symbol }}</td>
        <td>{{ asset.quantity }}</td>
        <td>{{ asset.price ? asset.price.toFixed(2) : '-' }}</td>
        <td>{{ asset.price ? (asset.quantity * asset.price).toFixed(2) : '-' }}</td>
          <td>
            <div class="action-buttons">
              <button class="edit-button" @click="openEditModal(asset)">Edit</button>
              <button class="delete-button" @click="deleteAsset(asset.id)">Delete</button>
            </div>
          </td>        
        </tr>
      </tbody>
    </table>
  <div v-if="showEditModal" class="modal">
    <div class="modal-content">
      <h2>Edit Asset</h2>
      <form @submit.prevent="editAsset">
        <div class="form-group">
          <label for="editSymbol">Asset Symbol:</label>
          <input type="text" id="editSymbol" v-model="editedAsset.symbol" required>
        </div>
        <div class="form-group">
          <label for="editQuantity">Quantity:</label>
          <input type="number" id="editQuantity" v-model.number="editedAsset.quantity" required>
        </div>
        <div class="form-group">
          <label for="editPrice">Current Market Price:</label>
          <input type="number" id="editPrice" v-model.number="editedAsset.price" step="0.01" required>
        </div>
        <div class="form-group">
          <button type="submit">Save</button>
          <button type="button" @click="closeEditModal">Cancel</button>
        </div>
      </form>
    </div>
  </div>
    <div class="add-asset-form">
      <h2>Add New Asset</h2>
      <form @submit.prevent="addAsset">
        <div class="form-group">
          <label for="symbol">Asset Symbol:</label>
          <input type="text" id="symbol" v-model="newAsset.symbol" required>
        </div>
        <div class="form-group">
          <label for="quantity">Quantity:</label>
          <input type="number" id="quantity" v-model.number="newAsset.quantity" required>
        </div>
        <div class="form-group">
          <label for="price">Current Market Price:</label>
          <input type="number" id="price" v-model.number="newAsset.price" step="0.01" required>
        </div>
        <div class="form-group">
          <button type="submit">Add Asset</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { Chart, registerables } from 'chart.js';
import { AssetServiceClient } from '../proto/asset_grpc_web_pb';
import { CreateAssetRequest, Empty, UpdateAssetRequest, DeleteAssetRequest } from '../proto/asset_pb';

const client = new AssetServiceClient('http://localhost:50052', null, null);
export default {
  data() {
    return {
      assets: [],
      newAsset: {
        symbol: '',
        quantity: null,
        price: null
      },
      chart: null,
      showEditModal: false,
      editedAsset: {
        id: '',
        symbol: '',
        quantity: null,
        price: null
      }
    };
  },
  mounted() {
    Chart.register(...registerables);
    this.renderDoughnutChart();
    this.listAssets();
    setInterval(() => {
      this.listAssets();
    }, 5000); // Update prices every 5 seconds
  },
  methods: {
    renderDoughnutChart() {
      if (this.chart) {
        this.chart.destroy();
      }
      const ctx = this.$refs.doughnutChart.getContext('2d');
      const labels = this.assets.map(asset => asset.symbol);
      const data = this.assets.map(asset => asset.quantity * asset.price);
      this.chart = new Chart(ctx, {
        type: 'doughnut',
        data: {
          labels: labels,
          datasets: [{
            data: data,
            backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40']
          }]
        },
        options: {
          responsive: true,
          title: {
            display: true,
            text: 'Assets Portfolio'
          }
        }
      });
    },
    openEditModal(asset) {
      this.editedAsset = { ...asset };
      this.showEditModal = true;
    },
    closeEditModal() {
      this.showEditModal = false;
    },
    editAsset() {
      const request = new UpdateAssetRequest();
      request.setId(this.editedAsset.id);
      request.setSymbol(this.editedAsset.symbol);
      request.setQuantity(this.editedAsset.quantity);
      request.setPrice(this.editedAsset.price);

      client.updateAsset(request, {}, (err, response) => {
        if (err) {
          console.error(err);
        } else {
          const updatedAsset = response.toObject();
          const index = this.assets.findIndex(a => a.id === updatedAsset.id);
          if (index !== -1) {
            this.assets.splice(index, 1, updatedAsset);
            this.renderDoughnutChart();
          }
          this.closeEditModal();
        }
      });
    },
    deleteAsset(assetId) {
      const request = new DeleteAssetRequest();
      request.setId(assetId);

      client.deleteAsset(request, {}, (err) => {
        if (err) {
          console.error(err);
        } else {
          this.assets = this.assets.filter(asset => asset.id !== assetId);
          this.renderDoughnutChart();
        }
      });
    },
    addAsset() {
      const request = new CreateAssetRequest();
      request.setSymbol(this.newAsset.symbol);
      request.setQuantity(this.newAsset.quantity);
      request.setPrice(this.newAsset.price);

      client.createAsset(request, {}, (err, response) => {
        if (err) {
          console.error(err);
        } else {
          const asset = response.toObject();
          this.assets.push(asset);
          this.newAsset = {
            symbol: '',
            quantity: null,
            price: null
          };
          this.renderDoughnutChart();
        }
      });
    },
    listAssets() {
      const request = new Empty();

      client.listAssets(request, {}, (err, response) => {
        if (err) {
          console.error(err);
        } else {
          const assetList = response.toObject().assetsList;
          this.assets = assetList.map(asset => ({
            id: asset.id,
            symbol: asset.symbol,
            quantity: asset.quantity,
            price: asset.price
          }));
          this.renderDoughnutChart();
        }
      });
    }
  },
};
</script>

<style>
body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
  margin: 0;
  padding: 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

@media screen and (max-width: 600px) {
  .container {
    padding: 10px;
  }

  h1 {
    font-size: 28px;
  }

  .add-asset-form {
    padding: 20px;
  }

  .form-group input {
    font-size: 14px;
  }

  .form-group button {
    font-size: 14px;
    padding: 10px 20px;
  }
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
  font-size: 36px;
  text-transform: uppercase;
  letter-spacing: 2px;
}

.chart-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 30px;
  background-color: #fff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

th, td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid #f2f2f2;
}

@media screen and (max-width: 600px) {
  table {
    font-size: 14px;
  }

  th, td {
    padding: 12px;
  }
}

th {
  background-color: #4CAF50;
  color: #fff;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 1px;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.action-buttons button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 1px;
  transition: background-color 0.3s ease;
}

.edit-button {
  background-color: #2196F3;
  color: #fff;
}

.edit-button:hover {
  background-color: #1976D2;
}

.delete-button {
  background-color: #f44336;
  color: #fff;
}

.delete-button:hover {
  background-color: #D32F2F;
}

.add-asset-form {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.add-asset-form h2 {
  margin-top: 0;
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: bold;
  margin-bottom: 10px;
  color: #555;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.form-group button {
  padding: 12px 24px;
  background-color: #4CAF50;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  text-transform: uppercase;
  letter-spacing: 1px;
  transition: background-color 0.3s ease;
}

.form-group button:hover {
  background-color: #45a049;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  max-width: 400px;
  width: 100%;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
  margin: 0;
  padding: 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

@media screen and (max-width: 600px) {
  .container {
    padding: 10px;
  }

  h1 {
    font-size: 28px;
  }

  .add-asset-form {
    padding: 20px;
  }

  .form-group input {
    font-size: 14px;
  }

  .form-group button {
    font-size: 14px;
    padding: 10px 20px;
  }
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
  font-size: 36px;
  text-transform: uppercase;
  letter-spacing: 2px;
}

.chart-container {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 30px;
  background-color: #fff;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

th, td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid #f2f2f2;
}

th {
  background-color: #4CAF50;
  color: #fff;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 1px;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.action-buttons button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 1px;
  transition: background-color 0.3s ease;
}

.edit-button {
  background-color: #2196F3;
  color: #fff;
}

.edit-button:hover {
  background-color: #1976D2;
}

.delete-button {
  background-color: #f44336;
  color: #fff;
}

.delete-button:hover {
  background-color: #D32F2F;
}

.add-asset-form {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.add-asset-form h2 {
  margin-top: 0;
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: bold;
  margin-bottom: 10px;
  color: #555;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.form-group button {
  padding: 12px 24px;
  background-color: #4CAF50;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  text-transform: uppercase;
  letter-spacing: 1px;
  transition: background-color 0.3s ease;
}

.form-group button:hover {
  background-color: #45a049;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  max-width: 400px;
  width: 100%;
}
</style>