<script>
  import { GetProduitParCodeBarres } from '../../wailsjs/go/main/App.js';

  let codeBarresSaisi = '';
  let panier = [];
  let messageErreur = '';
  let inputScanElement;

  $: totalHT = panier.reduce((acc, item) => acc + (item.prixVenteHt * item.quantite), 0);
  $: totalTTC = panier.reduce((acc, item) => acc + (item.prixVenteTtc * item.quantite), 0);
  $: totalTVA = totalTTC - totalHT;

  async function scannerOuValider(e) {
    if (e.key === 'Enter') {
      const code = codeBarresSaisi.trim();
      if (!code) return;
      messageErreur = '';

      try {
        const prod = await GetProduitParCodeBarres(code);
        ajouterAuPanier(prod);
        codeBarresSaisi = '';
      } catch (err) {
        messageErreur = 'Aucun produit trouvé pour le code : ' + code;
      }
    }
  }

  function ajouterAuPanier(p) {
    const index = panier.findIndex(item => item.id === p.id);
    if (index >= 0) {
      panier[index].quantite += 1;
      panier = [...panier];
    } else {
      panier = [
        ...panier,
        {
          id: p.id,
          codeBarres: p.codeBarres,
          reference: p.reference,
          designation: p.designation,
          prixVenteHt: p.prixVenteHt,
          prixVenteTtc: p.prixVenteTtc,
          tauxTva: p.tauxTva,
          quantite: 1,
          unite: p.unite
        }
      ];
    }
  }

  function modifierQuantite(index, delta) {
    panier[index].quantite += delta;
    if (panier[index].quantite <= 0) {
      supprimerLigne(index);
    } else {
      panier = [...panier];
    }
  }

  function supprimerLigne(index) {
    panier.splice(index, 1);
    panier = [...panier];
  }

  function viderPanier() {
    panier = [];
    messageErreur = '';
  }
</script>

<div class="caisse-container">
  <div class="header">
    <div>
      <h1>Caisse & Scan Rapide</h1>
      <p class="desc">Lecteur code-barres USB prêt à l'emploi (automatique code + Entrée)</p>
    </div>
    <div class="scan-bar">
      <input
        type="text"
        placeholder="Scannez un article (ou saisissez le code-barres + Entrée)..."
        bind:value={codeBarresSaisi}
        bind:this={inputScanElement}
        on:keydown={scannerOuValider}
      />
    </div>
  </div>

  {#if messageErreur}
    <div class="alert error">{messageErreur}</div>
  {/if}

  <div class="layout">
    <div class="panier-card">
      <div class="panier-header">
        <h2>Articles du panier ({panier.reduce((a, b) => a + b.quantite, 0)})</h2>
        {#if panier.length > 0}
          <button class="btn-vider" on:click={viderPanier}>Vider le panier</button>
        {/if}
      </div>

      {#if panier.length === 0}
        <div class="empty-state">
          <div class="empty-icon">🛒</div>
          <p>Le panier est vide.</p>
          <small>Scannez un code-barres avec votre douchette USB pour ajouter instantanément un produit.</small>
        </div>
      {:else}
        <table class="table">
          <thead>
            <tr>
              <th>Article</th>
              <th>Prix Unit. HT</th>
              <th>TVA</th>
              <th class="text-center">Quantité</th>
              <th>Total HT</th>
              <th>Total TTC</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {#each panier as item, idx}
              <tr>
                <td>
                  <strong>{item.designation}</strong>
                  <br /><small class="text-muted">{item.codeBarres}</small>
                </td>
                <td>{(item.prixVenteHt / 100).toFixed(2)} €</td>
                <td>{(item.tauxTva / 100)} %</td>
                <td class="text-center">
                  <div class="qty-control">
                    <button class="btn-qty" on:click={() => modifierQuantite(idx, -1)}>-</button>
                    <span class="qty-val">{item.quantite}</span>
                    <button class="btn-qty" on:click={() => modifierQuantite(idx, 1)}>+</button>
                  </div>
                </td>
                <td>{((item.prixVenteHt * item.quantite) / 100).toFixed(2)} €</td>
                <td><strong>{((item.prixVenteTtc * item.quantite) / 100).toFixed(2)} €</strong></td>
                <td class="text-right">
                  <button class="btn-remove" on:click={() => supprimerLigne(idx)}>&times;</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>

    <div class="reglement-card">
      <h2>Détail du Règlement</h2>
      <div class="recap-lines">
        <div class="recap-line">
          <span>Sous-total HT :</span>
          <span>{(totalHT / 100).toFixed(2)} €</span>
        </div>
        <div class="recap-line">
          <span>Montant TVA :</span>
          <span>{(totalTVA / 100).toFixed(2)} €</span>
        </div>
      </div>

      <div class="total-box">
        <span class="total-label">TOTAL TTC</span>
        <span class="total-val">{(totalTTC / 100).toFixed(2)} €</span>
      </div>

      <div class="mode-reglement">
        <label>Mode de paiement :</label>
        <div class="btn-group">
          <button class="btn btn-pay active">Espèces</button>
          <button class="btn btn-pay">Carte Bancaire</button>
          <button class="btn btn-pay">Virement</button>
        </div>
      </div>

      <button class="btn-valider" disabled={panier.length === 0}>
        Valider l'encaissement & Facturer
      </button>
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
    margin: 0 0 4px;
    font-size: 1.6rem;
    color: #0f172a;
  }
  .desc {
    margin: 0;
    color: #64748b;
  }
  .scan-bar input {
    width: 480px;
    padding: 12px 16px;
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
  .panier-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #f1f5f9;
    padding-bottom: 12px;
    margin-bottom: 16px;
  }
  .panier-header h2, .reglement-card h2 {
    margin: 0;
    font-size: 1.15rem;
    color: #1e293b;
  }
  .btn-vider {
    background: transparent;
    color: #dc2626;
    border: 1px solid #fca5a5;
    padding: 4px 10px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }
  .btn-vider:hover {
    background: #fee2e2;
  }
  .empty-state {
    text-align: center;
    padding: 80px 20px;
    color: #94a3b8;
  }
  .empty-icon {
    font-size: 3rem;
    margin-bottom: 10px;
  }
  .recap-lines {
    margin-bottom: 16px;
  }
  .recap-line {
    display: flex;
    justify-content: space-between;
    padding: 6px 0;
    color: #475569;
    font-size: 0.95rem;
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
    font-weight: 600;
  }
  .total-val {
    font-size: 2rem;
    font-weight: 800;
    color: #38bdf8;
  }
  .mode-reglement {
    margin-bottom: 20px;
  }
  .mode-reglement label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    color: #475569;
    margin-bottom: 6px;
  }
  .btn-group {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 8px;
  }
  .btn-pay {
    padding: 10px 6px;
    border: 1px solid #cbd5e1;
    background: #f8fafc;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.85rem;
    color: #334155;
  }
  .btn-pay.active {
    background: #0284c7;
    color: white;
    border-color: #0284c7;
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
  .btn-valider:hover:not(:disabled) {
    background: #15803d;
  }
  .btn-valider:disabled {
    background: #cbd5e1;
    cursor: not-allowed;
  }
  .table {
    width: 100%;
    border-collapse: collapse;
  }
  .table th, .table td {
    padding: 12px 10px;
    border-bottom: 1px solid #f1f5f9;
    text-align: left;
    font-size: 0.95rem;
  }
  .table th {
    background: #f8fafc;
    color: #64748b;
    font-size: 0.85rem;
  }
  .text-center { text-align: center; }
  .text-right { text-align: right; }
  .text-muted { color: #94a3b8; font-size: 0.8rem; }
  .qty-control {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: #f1f5f9;
    border-radius: 6px;
    padding: 2px 6px;
  }
  .btn-qty {
    background: white;
    border: 1px solid #cbd5e1;
    width: 26px;
    height: 26px;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
  }
  .qty-val {
    font-weight: 700;
    min-width: 20px;
  }
  .btn-remove {
    background: transparent;
    border: none;
    color: #94a3b8;
    font-size: 1.4rem;
    cursor: pointer;
  }
  .btn-remove:hover {
    color: #dc2626;
  }
  .alert {
    padding: 12px 16px;
    border-radius: 6px;
    margin-bottom: 16px;
    font-weight: 500;
  }
  .alert.error {
    background: #fef2f2;
    color: #991b1b;
    border: 1px solid #fecaca;
  }
</style>
