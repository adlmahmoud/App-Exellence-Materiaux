<script>
  import { onMount } from 'svelte';
  import { API } from '../lib/api.js';
  const { ListerProduits, CreerProduit, ModifierProduit, SupprimerProduit } = API;

  let produits = [];
  let total = 0;
  let page = 1;
  let totalPages = 1;
  let limit = 20;
  let recherche = '';
  let chargement = false;
  let messageErreur = '';
  let messageSucces = '';

  let modalNouveau = false;
  let modalEdition = false;
  let modalSuppression = false;
  let produitASupprimer = null;

  let formNouveau = {
    codeBarres: '',
    reference: '',
    designation: '',
    description: '',
    prixAchatHT: 0,
    prixVenteHT: 0,
    tauxTva: 2000,
    stockInitial: 0,
    stockMinimum: 5,
    unite: 'piece'
  };

  let formEdition = {
    id: '',
    codeBarres: '',
    reference: '',
    designation: '',
    description: '',
    prixAchatHT: 0,
    prixVenteHT: 0,
    tauxTva: 2000,
    stockMinimum: 5,
    unite: 'piece',
    actif: true
  };

  onMount(() => {
    chargerProduits();
  });

  async function chargerProduits() {
    chargement = true;
    messageErreur = '';
    try {
      const res = await ListerProduits({
        recherche: recherche.trim(),
        alerteStock: false,
        page: page,
        limit: limit,
        triChamp: 'designation',
        triAscendant: true
      });
      produits = res.produits || [];
      total = res.total || 0;
      page = res.page || 1;
      totalPages = res.totalPages || 1;
    } catch (err) {
      messageErreur = 'Erreur lors du chargement des produits : ' + err;
    } finally {
      chargement = false;
    }
  }

  let debounceTimer;

  function rechercheReactive() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      page = 1;
      chargerProduits();
    }, 150);
  }

  function effacerRecherche() {
    recherche = '';
    page = 1;
    chargerProduits();
  }

  function chercher() {
    clearTimeout(debounceTimer);
    page = 1;
    chargerProduits();
  }

  function ouvrirModalNouveau() {
    formNouveau = {
      codeBarres: '',
      reference: '',
      designation: '',
      description: '',
      prixAchatHT: 0,
      prixVenteHT: 0,
      tauxTva: 2000,
      stockInitial: 0,
      stockMinimum: 5,
      unite: 'piece'
    };
    messageErreur = '';
    messageSucces = '';
    modalNouveau = true;
  }

  async function soumettreNouveau() {
    messageErreur = '';
    try {
      const payload = {
        codeBarres: formNouveau.codeBarres.trim(),
        reference: formNouveau.reference.trim(),
        designation: formNouveau.designation.trim(),
        description: formNouveau.description.trim(),
        prixAchatHt: Math.round(Number(formNouveau.prixAchatHT) * 100),
        prixVenteHt: Math.round(Number(formNouveau.prixVenteHT) * 100),
        prixVenteTtc: 0,
        tauxTva: Number(formNouveau.tauxTva),
        stockInitial: Math.round(Number(formNouveau.stockInitial)),
        stockMinimum: Math.round(Number(formNouveau.stockMinimum)),
        unite: formNouveau.unite || 'piece'
      };

      await CreerProduit(payload);
      modalNouveau = false;
      messageSucces = 'Produit créé avec succès !';
      chargerProduits();
    } catch (err) {
      messageErreur = String(err);
    }
  }

  function ouvrirModalEdition(p) {
    formEdition = {
      id: p.id,
      codeBarres: p.codeBarres,
      reference: p.reference,
      designation: p.designation,
      description: p.description,
      prixAchatHT: (p.prixAchatHt / 100).toFixed(2),
      prixVenteHT: (p.prixVenteHt / 100).toFixed(2),
      tauxTva: p.tauxTva,
      stockMinimum: p.stockMinimum,
      unite: p.unite,
      actif: p.actif
    };
    modalEdition = true;
  }

  async function soumettreEdition() {
    messageErreur = '';
    try {
      const payload = {
        id: formEdition.id,
        codeBarres: formEdition.codeBarres.trim(),
        reference: formEdition.reference.trim(),
        designation: formEdition.designation.trim(),
        description: formEdition.description.trim(),
        prixAchatHt: Math.round(Number(formEdition.prixAchatHT) * 100),
        prixVenteHt: Math.round(Number(formEdition.prixVenteHT) * 100),
        prixVenteTtc: 0,
        tauxTva: Number(formEdition.tauxTva),
        stockMinimum: Math.round(Number(formEdition.stockMinimum)),
        unite: formEdition.unite,
        actif: formEdition.actif
      };

      await ModifierProduit(payload);
      modalEdition = false;
      messageSucces = 'Produit modifié avec succès !';
      chargerProduits();
    } catch (err) {
      messageErreur = String(err);
    }
  }

  function demanderSuppression(p) {
    produitASupprimer = p;
    modalSuppression = true;
  }

  async function confirmerSuppression() {
    if (!produitASupprimer) return;
    try {
      await SupprimerProduit(produitASupprimer.id);
      modalSuppression = false;
      produitASupprimer = null;
      messageSucces = 'Produit supprimé.';
      chargerProduits();
    } catch (err) {
      messageErreur = String(err);
    }
  }
</script>

<div class="page-container">
  <div class="header">
    <div>
      <h1>Catalogue Produits</h1>
      <p class="desc">{total} produit(s) enregistré(s) - Recherche instantanée FTS5</p>
    </div>
    <button class="btn-primary" on:click={ouvrirModalNouveau}>+ Nouveau Produit</button>
  </div>

  {#if messageSucces}
    <div class="alert success">{messageSucces}</div>
  {/if}
  {#if messageErreur}
    <div class="alert error">{messageErreur}</div>
  {/if}

  <div class="search-bar">
    <div class="search-input-wrapper">
      <span class="search-icon">🔍</span>
      <input
        type="text"
        placeholder="Rechercher instantanément par désignation, référence ou code-barres..."
        bind:value={recherche}
        on:input={rechercheReactive}
        on:keydown={(e) => e.key === 'Enter' && chercher()}
      />
      {#if recherche}
        <button class="btn-clear" on:click={effacerRecherche} title="Effacer la recherche">&times;</button>
      {/if}
    </div>
    <button class="btn-search" on:click={chercher}>Rechercher</button>
  </div>

  <div class="card">
    <table class="table">
      <thead>
        <tr>
          <th>Code-barres</th>
          <th>Référence</th>
          <th>Désignation</th>
          <th>Prix Achat HT</th>
          <th>Prix Vente HT</th>
          <th>TTC</th>
          <th>Stock</th>
          <th>Unité</th>
          <th class="text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if chargement}
          <tr><td colspan="9" class="empty">Chargement...</td></tr>
        {:else if produits.length === 0}
          <tr><td colspan="9" class="empty">Aucun produit trouvé. Cliquez sur "+ Nouveau Produit" pour commencer.</td></tr>
        {:else}
          {#each produits as p}
            <tr>
              <td><code>{p.codeBarres}</code></td>
              <td>{p.reference || '-'}</td>
              <td><strong>{p.designation}</strong></td>
              <td>{(p.prixAchatHt / 100).toFixed(2)} €</td>
              <td>{(p.prixVenteHt / 100).toFixed(2)} €</td>
              <td><span class="badge-ttc">{(p.prixVenteTtc / 100).toFixed(2)} €</span></td>
              <td>
                <span class="stock-pill {p.stockActuel <= p.stockMinimum ? 'stock-bas' : 'stock-ok'}">
                  {p.stockActuel}
                </span>
              </td>
              <td><span class="badge-unite">{p.unite}</span></td>
              <td class="text-right">
                <button class="btn-action edit" on:click={() => ouvrirModalEdition(p)}>Modifier</button>
                <button class="btn-action delete" on:click={() => demanderSuppression(p)}>Supprimer</button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>

    {#if totalPages > 1}
      <div class="pagination">
        <button disabled={page <= 1} on:click={() => { page--; chargerProduits(); }}>&laquo; Précédent</button>
        <span>Page {page} / {totalPages} ({total} articles)</span>
        <button disabled={page >= totalPages} on:click={() => { page++; chargerProduits(); }}>Suivant &raquo;</button>
      </div>
    {/if}
  </div>
</div>

{#if modalNouveau}
  <div class="overlay" on:click={() => modalNouveau = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Ajouter un Nouveau Produit</h2>
        <button class="close-btn" on:click={() => modalNouveau = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-row">
          <div class="form-group flex-2">
            <label>Code-barres (EAN / USB) * :</label>
            <input type="text" bind:value={formNouveau.codeBarres} placeholder="Ex: 3760123456789" required />
          </div>
          <div class="form-group flex-1">
            <label>Référence :</label>
            <input type="text" bind:value={formNouveau.reference} placeholder="Ex: MAT-001" />
          </div>
        </div>

        <div class="form-group">
          <label>Désignation de l'article * :</label>
          <input type="text" bind:value={formNouveau.designation} placeholder="Ex: Sac de ciment 35kg" required />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Prix Achat HT (€) :</label>
            <input type="number" step="0.01" bind:value={formNouveau.prixAchatHT} />
          </div>
          <div class="form-group">
            <label>Prix Vente HT (€) * :</label>
            <input type="number" step="0.01" bind:value={formNouveau.prixVenteHT} required />
          </div>
          <div class="form-group">
            <label>TVA :</label>
            <select bind:value={formNouveau.tauxTva}>
              <option value={2000}>20 %</option>
              <option value={1000}>10 %</option>
              <option value={550}>5.5 %</option>
              <option value={210}>2.1 %</option>
              <option value={0}>0 %</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Stock Initial :</label>
            <input type="number" bind:value={formNouveau.stockInitial} min="0" />
          </div>
          <div class="form-group">
            <label>Seuil d'alerte stock bas :</label>
            <input type="number" bind:value={formNouveau.stockMinimum} min="0" />
          </div>
          <div class="form-group">
            <label>Unité :</label>
            <input type="text" bind:value={formNouveau.unite} placeholder="piece, sac, m2..." />
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalNouveau = false}>Annuler</button>
        <button class="btn-primary" on:click={soumettreNouveau}>Enregistrer le produit</button>
      </div>
    </div>
  </div>
{/if}

{#if modalEdition}
  <div class="overlay" on:click={() => modalEdition = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Modifier le Produit</h2>
        <button class="close-btn" on:click={() => modalEdition = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-row">
          <div class="form-group flex-2">
            <label>Code-barres * :</label>
            <input type="text" bind:value={formEdition.codeBarres} required />
          </div>
          <div class="form-group flex-1">
            <label>Référence :</label>
            <input type="text" bind:value={formEdition.reference} />
          </div>
        </div>

        <div class="form-group">
          <label>Désignation * :</label>
          <input type="text" bind:value={formEdition.designation} required />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Prix Achat HT (€) :</label>
            <input type="number" step="0.01" bind:value={formEdition.prixAchatHT} />
          </div>
          <div class="form-group">
            <label>Prix Vente HT (€) * :</label>
            <input type="number" step="0.01" bind:value={formEdition.prixVenteHT} required />
          </div>
          <div class="form-group">
            <label>TVA :</label>
            <select bind:value={formEdition.tauxTva}>
              <option value={2000}>20 %</option>
              <option value={1000}>10 %</option>
              <option value={550}>5.5 %</option>
              <option value={210}>2.1 %</option>
              <option value={0}>0 %</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Seuil d'alerte stock :</label>
            <input type="number" bind:value={formEdition.stockMinimum} min="0" />
          </div>
          <div class="form-group">
            <label>Unité :</label>
            <input type="text" bind:value={formEdition.unite} />
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalEdition = false}>Annuler</button>
        <button class="btn-primary" on:click={soumettreEdition}>Mettre à jour</button>
      </div>
    </div>
  </div>
{/if}

{#if modalSuppression && produitASupprimer}
  <div class="overlay" on:click={() => modalSuppression = false}>
    <div class="modal modal-small" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Confirmation de suppression</h2>
        <button class="close-btn" on:click={() => modalSuppression = false}>&times;</button>
      </div>
      <div class="modal-body">
        <p>Êtes-vous sûr de vouloir supprimer définitivement le produit :</p>
        <p><strong>{produitASupprimer.designation}</strong> (<code>{produitASupprimer.codeBarres}</code>) ?</p>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalSuppression = false}>Annuler</button>
        <button class="btn-danger" on:click={confirmerSuppression}>Oui, supprimer</button>
      </div>
    </div>
  </div>
{/if}

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
  .search-bar {
    display: flex;
    gap: 10px;
    margin-bottom: 16px;
  }
  .search-input-wrapper {
    flex: 1;
    position: relative;
    display: flex;
    align-items: center;
  }
  .search-icon {
    position: absolute;
    left: 12px;
    font-size: 1rem;
    color: #64748b;
    pointer-events: none;
  }
  .search-input-wrapper input {
    width: 100%;
    padding: 10px 38px 10px 38px;
    border: 1px solid #cbd5e1;
    border-radius: 8px;
    font-size: 1rem;
    outline: none;
    box-sizing: border-box;
    transition: border-color 0.15s, box-shadow 0.15s;
  }
  .search-input-wrapper input:focus {
    border-color: #0284c7;
    box-shadow: 0 0 0 3px rgba(2, 132, 199, 0.15);
  }
  .btn-clear {
    position: absolute;
    right: 10px;
    background: transparent;
    border: none;
    font-size: 1.3rem;
    color: #94a3b8;
    cursor: pointer;
    line-height: 1;
    padding: 4px 6px;
  }
  .btn-clear:hover {
    color: #334155;
  }
  .btn-search {
    background: #334155;
    color: white;
    border: none;
    padding: 0 20px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-primary {
    background: #0284c7;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-secondary {
    background: #f1f5f9;
    color: #475569;
    border: 1px solid #cbd5e1;
    padding: 10px 18px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-danger {
    background: #dc2626;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
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
    font-size: 0.95rem;
  }
  .table th {
    background: #f8fafc;
    color: #475569;
    font-size: 0.85rem;
    font-weight: 600;
  }
  .text-right {
    text-align: right;
  }
  .empty {
    text-align: center;
    color: #94a3b8;
    padding: 40px !important;
  }
  .stock-pill {
    padding: 4px 10px;
    border-radius: 12px;
    font-weight: 700;
    font-size: 0.85rem;
  }
  .stock-ok {
    background: #dcfce7;
    color: #166534;
  }
  .stock-bas {
    background: #fee2e2;
    color: #991b1b;
  }
  .badge-ttc {
    font-weight: 700;
    color: #0284c7;
  }
  .badge-unite {
    background: #f1f5f9;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    color: #64748b;
  }
  .btn-action {
    background: transparent;
    border: 1px solid #cbd5e1;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.85rem;
    cursor: pointer;
    margin-left: 4px;
  }
  .btn-action.edit:hover {
    background: #0284c7;
    color: white;
    border-color: #0284c7;
  }
  .btn-action.delete:hover {
    background: #dc2626;
    color: white;
    border-color: #dc2626;
  }
  .pagination {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    background: #f8fafc;
    border-top: 1px solid #e2e8f0;
  }
  .pagination button {
    background: #ffffff;
    border: 1px solid #cbd5e1;
    padding: 6px 14px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 500;
  }
  .pagination button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .overlay {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(15, 23, 42, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .modal {
    background: white;
    border-radius: 10px;
    width: 600px;
    max-width: 90vw;
    box-shadow: 0 20px 25px -5px rgba(0,0,0,0.2);
    overflow: hidden;
  }
  .modal-small {
    width: 440px;
  }
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    background: #f8fafc;
    border-bottom: 1px solid #e2e8f0;
  }
  .modal-header h2 {
    margin: 0;
    font-size: 1.15rem;
    color: #0f172a;
  }
  .close-btn {
    background: transparent;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #64748b;
  }
  .modal-body {
    padding: 20px;
  }
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    padding: 14px 20px;
    background: #f8fafc;
    border-top: 1px solid #e2e8f0;
  }
  .form-group {
    margin-bottom: 14px;
  }
  .form-row {
    display: flex;
    gap: 12px;
  }
  .form-row .form-group {
    flex: 1;
  }
  .flex-2 { flex: 2 !important; }
  .flex-1 { flex: 1 !important; }
  label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    color: #475569;
    margin-bottom: 4px;
  }
  input, select {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 0.95rem;
    box-sizing: border-box;
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
