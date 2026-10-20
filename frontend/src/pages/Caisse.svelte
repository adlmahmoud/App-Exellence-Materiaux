<script>
  let codeBarres = '';
  let panier = [];
  let totalTTC = 0;

  function onScan(e) {
    if (e.key === 'Enter' && codeBarres.trim() !== '') {
      codeBarres = '';
    }
  }
</script>

<div class="caisse-container">
  <div class="header">
    <h1>Caisse & Vente Rapide</h1>
    <div class="scan-bar">
      <input
        type="text"
        placeholder="Scannez un article (ou saisissez le code-barres + Entrée)..."
        bind:value={codeBarres}
        on:keydown={onScan}
        autofocus
      />
    </div>
  </div>

  <div class="layout">
    <div class="panier-card">
      <h2>Articles en cours</h2>
      {#if panier.length === 0}
        <div class="empty-state">
          <p>Le panier est vide.</p>
          <small>Scannez un code-barres avec le lecteur USB pour ajouter un produit.</small>
        </div>
      {:else}
        <table class="table">
          <thead>
            <tr>
              <th>Désignation</th>
              <th>Qté</th>
              <th>PU HT</th>
              <th>TVA</th>
              <th>Total HT</th>
            </tr>
          </thead>
          <tbody>
            {#each panier as item}
              <tr>
                <td>{item.designation}</td>
                <td>{item.quantite}</td>
                <td>{(item.prixUnitaireHt / 100).toFixed(2)} €</td>
                <td>{(item.tauxTva / 100)} %</td>
                <td>{(item.montantHt / 100).toFixed(2)} €</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>

    <div class="reglement-card">
      <h2>Total à régler</h2>
      <div class="total-box">
        <span class="total-label">TOTAL TTC</span>
        <span class="total-val">{(totalTTC / 100).toFixed(2)} €</span>
      </div>
      <div class="btn-group">
        <button class="btn btn-especes">Espèces</button>
        <button class="btn btn-cb">Carte Bancaire</button>
        <button class="btn btn-virement">Virement</button>
      </div>
      <button class="btn-valider">Valider la vente & Facturer</button>
    </div>
  </div>
</div>

<style>
  .caisse-container {
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
    margin: 0;
    font-size: 1.6rem;
    color: #0f172a;
  }
  .scan-bar input {
    width: 480px;
    padding: 10px 16px;
    border: 2px solid #0284c7;
    border-radius: 8px;
    font-size: 1rem;
    outline: none;
    box-shadow: 0 0 0 3px rgba(2, 132, 199, 0.15);
  }
  .layout {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 20px;
  }
  .panier-card, .reglement-card {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
    padding: 20px;
  }
  .panier-card h2, .reglement-card h2 {
    margin: 0 0 16px;
    font-size: 1.1rem;
    color: #1e293b;
    border-bottom: 1px solid #f1f5f9;
    padding-bottom: 8px;
  }
  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: #94a3b8;
  }
  .total-box {
    background: #0f172a;
    color: white;
    padding: 20px;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  .total-label {
    font-size: 1rem;
    font-weight: 500;
  }
  .total-val {
    font-size: 1.8rem;
    font-weight: bold;
    color: #38bdf8;
  }
  .btn-group {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 8px;
    margin-bottom: 16px;
  }
  .btn {
    padding: 12px 8px;
    border: 1px solid #cbd5e1;
    background: #f8fafc;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
    color: #334155;
  }
  .btn:hover {
    background: #e2e8f0;
  }
  .btn-valider {
    width: 100%;
    padding: 14px;
    background: #16a34a;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1.1rem;
    font-weight: bold;
    cursor: pointer;
  }
  .btn-valider:hover {
    background: #15803d;
  }
  .table {
    width: 100%;
    border-collapse: collapse;
  }
  .table th, .table td {
    padding: 10px;
    border-bottom: 1px solid #f1f5f9;
    text-align: left;
  }
  .table th {
    background: #f8fafc;
    color: #64748b;
    font-size: 0.85rem;
  }
</style>
