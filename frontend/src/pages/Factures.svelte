<script>
</script>

<div class="page-container">
  <div class="header">
    <div>
      <h1>Factures & Bons</h1>
      <p class="desc">Création, numérotation séquentielle sans trou et impression PDF</p>
    </div>
    <div class="actions">
      <button class="btn btn-primary">+ Nouvelle Facture</button>
      <button class="btn btn-secondary">+ Nouveau Bon</button>
    </div>
  </div>

  <div class="card">
    <table class="table">
      <thead>
        <tr>
          <th>Numéro</th>
          <th>Date</th>
          <th>Client</th>
          <th>Total HT</th>
          <th>Total TTC</th>
          <th>Statut</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td colspan="7" class="empty">Aucune facture enregistrée.</td>
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
  .actions {
    display: flex;
    gap: 10px;
  }
  .btn {
    padding: 10px 16px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    border: none;
  }
  .btn-primary {
    background: #0284c7;
    color: white;
  }
  .btn-secondary {
    background: #475569;
    color: white;
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
