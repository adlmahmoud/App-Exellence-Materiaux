<script>
  import {
    GetProduitParCodeBarres,
    ListerProduits,
    CreerFacture,
    ImprimerFacturePDF,
    CreerBon,
    ImprimerBonPDF
  } from '../../wailsjs/go/main/App.js';

  let codeBarresSaisi = '';
  let panier = [];
  let messageErreur = '';
  let messageSucces = '';
  let inputScanElement;

  let suggestions = [];
  let timerSuggestions;

  let clientNom = 'Client Comptoir';
  let modeReglement = 'especes';
  let chargementPaiement = false;

  $: totalHT = panier.reduce((acc, item) => acc + (item.prixVenteHt * item.quantite), 0);
  $: totalTTC = panier.reduce((acc, item) => acc + (item.prixVenteTtc * item.quantite), 0);
  $: totalTVA = totalTTC - totalHT;

  function onInputScan() {
    clearTimeout(timerSuggestions);
    const q = codeBarresSaisi.trim();
    if (q.length < 2) {
      suggestions = [];
      return;
    }
    timerSuggestions = setTimeout(async () => {
      try {
        const res = await ListerProduits({ recherche: q, limit: 5 });
        suggestions = res.produits || [];
      } catch (e) {
        suggestions = [];
      }
    }, 150);
  }

  function choisirSuggestion(p) {
    ajouterAuPanier(p);
    codeBarresSaisi = '';
    suggestions = [];
    if (inputScanElement) inputScanElement.focus();
  }

  async function scannerOuValider(e) {
    if (e.key === 'Enter') {
      const code = codeBarresSaisi.trim();
      if (!code) return;
      messageErreur = '';
      messageSucces = '';
      suggestions = [];

      try {
        const prod = await GetProduitParCodeBarres(code);
        ajouterAuPanier(prod);
        codeBarresSaisi = '';
      } catch (err) {
        try {
          const res = await ListerProduits({ recherche: code, limit: 1 });
          if (res && res.produits && res.produits.length > 0) {
            ajouterAuPanier(res.produits[0]);
            codeBarresSaisi = '';
          } else {
            messageErreur = 'Aucun produit trouvé pour : ' + code;
          }
        } catch (e2) {
          messageErreur = 'Aucun produit trouvé pour : ' + code;
        }
      }
    } else if (e.key === 'Escape') {
      suggestions = [];
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
    messageSucces = '';
  }

  async function validerEtFacturer() {
    if (panier.length === 0) return;
    chargementPaiement = true;
    messageErreur = '';
    messageSucces = '';

    try {
      const lignes = panier.map(item => ({
        produitId: item.id,
        codeBarres: item.codeBarres || '',
        reference: item.reference || '',
        designation: item.designation,
        quantite: Math.round(Number(item.quantite)),
        prixUnitaireHt: item.prixVenteHt,
        tauxTva: item.tauxTva
      }));

      const res = await CreerFacture({
        clientNom: clientNom.trim() || 'Client Comptoir',
        modeReglement: modeReglement,
        lignes: lignes
      });

      messageSucces = `Encaissement réussi ! Facture ${res.facture.numero} générée et ouverte en PDF.`;
      panier = [];

      try {
        await ImprimerFacturePDF(res.facture.id);
      } catch (errPdf) {}

      if (inputScanElement) inputScanElement.focus();
    } catch (err) {
      messageErreur = 'Erreur lors de la validation : ' + err;
    } finally {
      chargementPaiement = false;
    }
  }

  async function validerEtBonLivraison() {
    if (panier.length === 0) return;
    chargementPaiement = true;
    messageErreur = '';
    messageSucces = '';

    try {
      const lignes = panier.map(item => ({
        produitId: item.id,
        codeBarres: item.codeBarres || '',
        reference: item.reference || '',
        designation: item.designation,
        quantite: Math.round(Number(item.quantite)),
        prixUnitaireHt: item.prixVenteHt,
        tauxTva: item.tauxTva
      }));

      const res = await CreerBon({
        typeBon: 'livraison',
        clientNom: clientNom.trim() || 'Client Chantier',
        lignes: lignes
      });

      messageSucces = `Bon de livraison ${res.bon.numero} généré et ouvert en PDF.`;
      panier = [];

      try {
        await ImprimerBonPDF(res.bon.id);
      } catch (errPdf) {}

      if (inputScanElement) inputScanElement.focus();
    } catch (err) {
      messageErreur = 'Erreur lors de la création du bon : ' + err;
    } finally {
      chargementPaiement = false;
    }
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
        placeholder="Scannez ou tapez (code-barres, nom d'article, référence)..."
        bind:value={codeBarresSaisi}
        bind:this={inputScanElement}
        on:input={onInputScan}
        on:keydown={scannerOuValider}
      />
      {#if suggestions.length > 0}
        <div class="suggestions-dropdown">
          {#each suggestions as s}
            <div class="suggestion-item" on:click={() => choisirSuggestion(s)}>
              <span class="sugg-name">{s.designation}</span>
              <span class="sugg-code">{s.codeBarres}</span>
              <span class="sugg-price">{(s.prixVenteTtc / 100).toFixed(2)} € TTC</span>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>

  {#if messageSucces}
    <div class="alert success">{messageSucces}</div>
  {/if}
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

      <div class="client-box">
        <label>Client :</label>
        <input type="text" bind:value={clientNom} placeholder="Nom du client" />
      </div>

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
          <button
            class="btn btn-pay {modeReglement === 'especes' ? 'active' : ''}"
            on:click={() => modeReglement = 'especes'}
          >
            Espèces
          </button>
          <button
            class="btn btn-pay {modeReglement === 'cb' ? 'active' : ''}"
            on:click={() => modeReglement = 'cb'}
          >
            Carte Bancaire
          </button>
          <button
            class="btn btn-pay {modeReglement === 'virement' ? 'active' : ''}"
            on:click={() => modeReglement = 'virement'}
          >
            Virement
          </button>
          <button
            class="btn btn-pay {modeReglement === 'cheque' ? 'active' : ''}"
            on:click={() => modeReglement = 'cheque'}
          >
            Chèque
          </button>
        </div>
      </div>

      <button
        class="btn-valider"
        disabled={panier.length === 0 || chargementPaiement}
        on:click={validerEtFacturer}
      >
        {chargementPaiement ? 'Validation en cours...' : 'Valider l\'encaissement & Facturer (PDF)'}
      </button>

      <button
        class="btn-bon"
        disabled={panier.length === 0 || chargementPaiement}
        on:click={validerEtBonLivraison}
      >
        Générer Bon de livraison
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
  .scan-bar {
    position: relative;
    width: 500px;
  }
  .scan-bar input {
    width: 100%;
    padding: 12px 16px;
    border: 2px solid #0284c7;
    border-radius: 8px;
    font-size: 1rem;
    outline: none;
    box-shadow: 0 0 0 3px rgba(2, 132, 199, 0.15);
    box-sizing: border-box;
  }
  .suggestions-dropdown {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    background: #ffffff;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    z-index: 100;
    max-height: 260px;
    overflow-y: auto;
  }
  .suggestion-item {
    padding: 10px 14px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    border-bottom: 1px solid #f1f5f9;
    font-size: 0.95rem;
    transition: background 0.1s;
  }
  .suggestion-item:last-child {
    border-bottom: none;
  }
  .suggestion-item:hover {
    background: #f0f9ff;
  }
  .sugg-name {
    font-weight: 600;
    color: #0f172a;
    flex: 1;
  }
  .sugg-code {
    font-size: 0.8rem;
    color: #64748b;
    margin: 0 12px;
  }
  .sugg-price {
    font-weight: 700;
    color: #0284c7;
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
    box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  }
  .panier-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    border-bottom: 1px solid #f1f5f9;
    padding-bottom: 12px;
  }
  .panier-header h2 {
    margin: 0;
    font-size: 1.2rem;
    color: #1e293b;
  }
  .btn-vider {
    background: transparent;
    color: #dc2626;
    border: 1px solid #fca5a5;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
  }
  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: #64748b;
  }
  .empty-icon {
    font-size: 3rem;
    margin-bottom: 12px;
  }
  .empty-state p {
    font-size: 1.1rem;
    font-weight: 600;
    margin: 0 0 6px;
    color: #1e293b;
  }
  .table {
    width: 100%;
    border-collapse: collapse;
  }
  .table th, .table td {
    padding: 10px 12px;
    border-bottom: 1px solid #f1f5f9;
    text-align: left;
    font-size: 0.95rem;
  }
  .table th {
    background: #f8fafc;
    color: #475569;
    font-size: 0.85rem;
    font-weight: 600;
  }
  .text-center {
    text-align: center;
  }
  .text-right {
    text-align: right;
  }
  .text-muted {
    color: #94a3b8;
  }
  .qty-control {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: #f1f5f9;
    border-radius: 4px;
    padding: 2px 4px;
  }
  .btn-qty {
    background: #ffffff;
    border: 1px solid #cbd5e1;
    width: 24px;
    height: 24px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-weight: bold;
    font-size: 0.9rem;
  }
  .qty-val {
    font-weight: 700;
    min-width: 20px;
    text-align: center;
  }
  .btn-remove {
    background: transparent;
    border: none;
    color: #ef4444;
    font-size: 1.2rem;
    cursor: pointer;
  }
  .reglement-card h2 {
    margin: 0 0 16px;
    font-size: 1.2rem;
    color: #1e293b;
    border-bottom: 1px solid #f1f5f9;
    padding-bottom: 10px;
  }
  .client-box {
    margin-bottom: 16px;
  }
  .client-box label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    color: #475569;
    margin-bottom: 4px;
  }
  .client-box input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 0.95rem;
    box-sizing: border-box;
  }
  .recap-lines {
    margin-bottom: 16px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .recap-line {
    display: flex;
    justify-content: space-between;
    font-size: 0.95rem;
    color: #475569;
  }
  .total-box {
    background: #0f172a;
    color: white;
    padding: 16px;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  .total-label {
    font-size: 1rem;
    font-weight: 600;
    color: #94a3b8;
  }
  .total-val {
    font-size: 1.6rem;
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
    margin-bottom: 8px;
  }
  .btn-group {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }
  .btn-pay {
    background: #f8fafc;
    border: 1px solid #cbd5e1;
    padding: 10px;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    color: #475569;
    cursor: pointer;
    transition: all 0.15s;
  }
  .btn-pay.active {
    background: #0284c7;
    color: white;
    border-color: #0284c7;
  }
  .btn-valider {
    width: 100%;
    background: #16a34a;
    color: white;
    border: none;
    padding: 14px;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 700;
    cursor: pointer;
    margin-bottom: 10px;
    transition: background 0.15s;
  }
  .btn-valider:hover:not(:disabled) {
    background: #15803d;
  }
  .btn-valider:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .btn-bon {
    width: 100%;
    background: #475569;
    color: white;
    border: none;
    padding: 10px;
    border-radius: 8px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-bon:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .alert {
    padding: 12px 16px;
    border-radius: 6px;
    margin-bottom: 16px;
    font-weight: 500;
  }
  .alert.success {
    background: #ecfdf5;
    color: #065f46;
    border: 1px solid #a7f3d0;
  }
  .alert.error {
    background: #fef2f2;
    color: #991b1b;
    border: 1px solid #fecaca;
  }
</style>
