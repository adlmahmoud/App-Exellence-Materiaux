<script>
  let recherche = '';
</script>

<div class="page-container">
  <div class="header">
    <div>
      <h1>Catalogue Produits</h1>
      <p class="desc">Gestion des articles, prix au centime et codes-barres</p>
    </div>
    <button class="btn-ajouter">+ Nouveau Produit</button>
  </div>

  <div class="search-bar">
    <input
      type="text"
      placeholder="Recherche rapide par nom, référence ou code-barres (FTS5)..."
      bind:value={recherche}
    />
  </div>

  <div class="card">
    <table class="table">
      <thead>
        <tr>
          <th>Code-barres</th>
          <th>Référence</th>
          <th>Désignation</th>
          <th>Prix Vente HT</th>
          <th>TVA</th>
          <th>Stock</th>
          <th>Unité</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td colspan="8" class="empty">Aucun produit pour le moment. Cliquez sur "+ Nouveau Produit" pour en ajouter.</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>

<style>
  .page-container {
    padding: 24px;
    max-width: 1200px;
    margin: 0 auto;
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  .header h1 {
    margin: 0 0 4px;
    font-size: 1.6rem;
    color: #0f172a;
  }
  .desc {
    margin: 0;
    color: #64748b;
  }
  .btn-ajouter {
    background: #0284c7;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .search-bar {
    margin-bottom: 16px;
  }
  .search-bar input {
    width: 100%;
    padding: 10px 16px;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    font-size: 1rem;
    box-sizing: border-box;
  }
  .card {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
    overflow: hidden;
  }
  .table {
    width: 100%;
    border-collapse: collapse;
  }
  .table th, .table td {
    padding: 12px 16px;
    border-bottom: 1px solid #f1f5f9;
    text-align: left;
  }
  .table th {
    background: #f8fafc;
    color: #475569;
    font-size: 0.85rem;
    font-weight: 600;
  }
  .empty {
    text-align: center;
    color: #94a3b8;
    padding: 40px !important;
  }
</style>
