<script>
  import { onMount } from 'svelte';
  import {
    ListerMouvementsStock,
    CreerMouvementStock,
    SupprimerMouvementStock,
    AjusterStock,
    ObtenirAlertesStock,
    ListerProduits
  } from '../lib/api.js';

  let mouvements = [];
  let alertes = [];
  let produitsListe = [];
  let total = 0;
  let page = 1;
  let totalPages = 1;
  let limit = 20;
  let chargement = false;
  let messageSucces = '';
  let messageErreur = '';

  let modalEntree = false;
  let modalSortie = false;
  let modalAjustement = false;
  let mouvementASupprimer = null;

  let formEntree = {
    produitId: '',
    quantite: 1,
    motif: 'Réception fournisseur',
    referenceDocument: ''
  };

  let formSortie = {
    produitId: '',
    quantite: 1,
    motif: 'Casse / Perte',
    referenceDocument: ''
  };

  let formAjustement = {
    produitId: '',
    nouveauStock: 0,
    motif: 'Inventaire physique'
  };

  onMount(() => {
    chargerDonnees();
    chargerProduitsPourSelect();
  });

  async function chargerDonnees() {
    chargement = true;
    messageErreur = '';
    try {
      const resMouvements = await ListerMouvementsStock({
        page: page,
        limit: limit
      });
      mouvements = resMouvements.mouvements || [];
      total = resMouvements.total || 0;
      page = resMouvements.page || 1;
      totalPages = resMouvements.totalPages || 1;

      const resAlertes = await ObtenirAlertesStock();
      alertes = resAlertes || [];
    } catch (err) {
      messageErreur = 'Erreur chargement stock : ' + err;
    } finally {
      chargement = false;
    }
  }

  async function chargerProduitsPourSelect() {
    try {
      const res = await ListerProduits({ limit: 100, page: 1, triChamp: 'designation', triAscendant: true });
      produitsListe = res.produits || [];
      if (produitsListe.length > 0) {
        if (!formEntree.produitId) formEntree.produitId = produitsListe[0].id;
        if (!formSortie.produitId) formSortie.produitId = produitsListe[0].id;
        if (!formAjustement.produitId) {
          formAjustement.produitId = produitsListe[0].id;
          formAjustement.nouveauStock = produitsListe[0].stockActuel;
        }
      }
    } catch (e) {}
  }

  function ouvrirModalEntree() {
    messageSucces = '';
    messageErreur = '';
    if (produitsListe.length > 0 && !formEntree.produitId) {
      formEntree.produitId = produitsListe[0].id;
    }
    formEntree.quantite = 1;
    formEntree.motif = 'Réception fournisseur';
    formEntree.referenceDocument = '';
    modalEntree = true;
  }

  async function validerEntree() {
    messageErreur = '';
    if (!formEntree.produitId) {
      messageErreur = 'Veuillez sélectionner un produit.';
      return;
    }
    if (Number(formEntree.quantite) <= 0) {
      messageErreur = 'La quantité doit être supérieure à zéro.';
      return;
    }

    try {
      await CreerMouvementStock({
        produitId: formEntree.produitId,
        typeMouvement: 'entree',
        quantite: Math.round(Number(formEntree.quantite)),
        motif: formEntree.motif.trim(),
        referenceDocument: formEntree.referenceDocument.trim()
      });
      modalEntree = false;
      messageSucces = 'Entrée de stock enregistrée avec succès !';
      chargerDonnees();
      chargerProduitsPourSelect();
    } catch (err) {
      messageErreur = String(err);
    }
  }

  function ouvrirModalSortie() {
    messageSucces = '';
    messageErreur = '';
    if (produitsListe.length > 0 && !formSortie.produitId) {
      formSortie.produitId = produitsListe[0].id;
    }
    formSortie.quantite = 1;
    formSortie.motif = 'Casse / Perte';
    formSortie.referenceDocument = '';
    modalSortie = true;
  }

  async function validerSortie() {
    messageErreur = '';
    if (!formSortie.produitId) {
      messageErreur = 'Veuillez sélectionner un produit.';
      return;
    }
    if (Number(formSortie.quantite) <= 0) {
      messageErreur = 'La quantité doit être supérieure à zéro.';
      return;
    }

    try {
      await CreerMouvementStock({
        produitId: formSortie.produitId,
        typeMouvement: 'sortie',
        quantite: Math.round(Number(formSortie.quantite)),
        motif: formSortie.motif.trim(),
        referenceDocument: formSortie.referenceDocument.trim()
      });
      modalSortie = false;
      messageSucces = 'Sortie de stock enregistrée !';
      chargerDonnees();
      chargerProduitsPourSelect();
    } catch (err) {
      messageErreur = String(err);
    }
  }

  function ouvrirModalAjustement() {
    messageSucces = '';
    messageErreur = '';
    if (produitsListe.length > 0) {
      formAjustement.produitId = produitsListe[0].id;
      formAjustement.nouveauStock = produitsListe[0].stockActuel;
    }
    formAjustement.motif = 'Inventaire physique';
    modalAjustement = true;
  }

  function onChangementProduitAjustement(e) {
    const p = produitsListe.find(item => item.id === e.target.value);
    if (p) {
      formAjustement.nouveauStock = p.stockActuel;
    }
  }

  async function validerAjustement() {
    messageErreur = '';
    if (!formAjustement.produitId) {
      messageErreur = 'Veuillez sélectionner un produit.';
      return;
    }

    try {
      await AjusterStock({
        produitId: formAjustement.produitId,
        nouveauStock: Math.round(Number(formAjustement.nouveauStock)),
        motif: formAjustement.motif.trim()
      });
      modalAjustement = false;
      messageSucces = 'Stock ajusté avec succès !';
      chargerDonnees();
      chargerProduitsPourSelect();
    } catch (err) {
      messageErreur = String(err);
    }
  }

  async function confirmerSuppressionMouvement() {
    if (!mouvementASupprimer) return;
    messageErreur = '';
    try {
      await SupprimerMouvementStock(mouvementASupprimer.id);
      messageSucces = 'Mouvement supprimé de l’historique. Le stock actuel n’a pas été modifié.';
      if (mouvements.length === 1 && page > 1) page--;
      mouvementASupprimer = null;
      await chargerDonnees();
    } catch (err) {
      messageErreur = String(err);
    }
  }

  function formaterDate(dStr) {
    if (!dStr) return '-';
    try {
      const d = new Date(dStr);
      return d.toLocaleDateString('fr-FR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (e) {
      return dStr;
    }
  }
</script>

<div class="page-container">
  <div class="header">
    <div>
      <h1>Gestion des Stocks</h1>
      <p class="desc">Entrées, sorties, ajustements d'inventaire et alertes</p>
    </div>
    <div class="actions">
      <button class="btn btn-entree" on:click={ouvrirModalEntree}>+ Entrée de stock</button>
      <button class="btn btn-sortie" on:click={ouvrirModalSortie}>- Sortie de stock</button>
      <button class="btn btn-ajust" on:click={ouvrirModalAjustement}>Ajustement d'inventaire</button>
    </div>
  </div>

  {#if messageSucces}
    <div class="alert success">{messageSucces}</div>
  {/if}
  {#if messageErreur}
    <div class="alert error">{messageErreur}</div>
  {/if}

  {#if alertes.length > 0}
    <div class="alert-box">
      <div class="alert-header">
        <span class="warning-icon">⚠️</span>
        <strong>Alertes Stock Bas ({alertes.length} article(s) sous le seuil critique)</strong>
      </div>
      <div class="alert-list">
        {#each alertes as a}
          <div class="alert-pill">
            <span class="name">{a.designation}</span>
            <span class="stock-info">Stock actuel : <strong>{a.stockActuel}</strong> (Seuil : {a.stockMinimum})</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <div class="card">
    <div class="card-header">
      <h2>Historique des Mouvements ({total})</h2>
      <button class="btn-refresh" on:click={chargerDonnees}>Actualiser</button>
    </div>

    <table class="table">
      <thead>
        <tr>
          <th>Date & Heure</th>
          <th>Type</th>
          <th>Article</th>
          <th class="text-center">Quantité</th>
          <th class="text-center">Avant</th>
          <th class="text-center">Après</th>
          <th>Motif / Document</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        {#if chargement}
          <tr><td colspan="8" class="empty">Chargement de l'historique...</td></tr>
        {:else if mouvements.length === 0}
          <tr><td colspan="8" class="empty">Aucun mouvement de stock enregistré.</td></tr>
        {:else}
          {#each mouvements as m}
            <tr>
              <td>{formaterDate(m.createdAt)}</td>
              <td>
                {#if m.typeMouvement === 'entree'}
                  <span class="badge badge-entree">+ Entrée</span>
                {:else if m.typeMouvement === 'sortie'}
                  <span class="badge badge-sortie">- Sortie</span>
                {:else if m.typeMouvement === 'vente'}
                  <span class="badge badge-vente">Vente</span>
                {:else}
                  <span class="badge badge-ajust">Ajustement</span>
                {/if}
              </td>
              <td>
                <strong>{m.produitDesignation}</strong>
                {#if m.produitCodeBarres}
                  <br /><small class="text-muted">{m.produitCodeBarres}</small>
                {/if}
              </td>
              <td class="text-center font-bold">
                {#if m.typeMouvement === 'entree'}
                  <span class="text-green">+{m.quantite}</span>
                {:else if m.typeMouvement === 'sortie' || m.typeMouvement === 'vente'}
                  <span class="text-red">-{m.quantite}</span>
                {:else}
                  <span>{m.quantite > 0 ? '+' : ''}{m.quantite}</span>
                {/if}
              </td>
              <td class="text-center">{m.stockAvant}</td>
              <td class="text-center font-bold">{m.stockApres}</td>
              <td>
                <span>{m.motif || '-'}</span>
                {#if m.referenceDocument}
                  <span class="doc-ref">({m.referenceDocument})</span>
                {/if}
              </td>
              <td>
                <button class="btn-delete-movement" on:click={() => mouvementASupprimer = m}>
                  Supprimer
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>

    {#if totalPages > 1}
      <div class="pagination">
        <button disabled={page <= 1} on:click={() => { page--; chargerDonnees(); }}>&laquo; Précédent</button>
        <span>Page {page} / {totalPages} ({total} mouvements)</span>
        <button disabled={page >= totalPages} on:click={() => { page++; chargerDonnees(); }}>Suivant &raquo;</button>
      </div>
    {/if}
  </div>
</div>

{#if mouvementASupprimer}
  <div class="overlay">
    <div class="modal modal-small" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Confirmation de suppression</h2>
        <button class="close-btn" on:click={() => mouvementASupprimer = null}>&times;</button>
      </div>
      <div class="modal-body">
        <p>Supprimer définitivement ce mouvement de l’historique ?</p>
        <p><strong>{mouvementASupprimer.designation}</strong> — {mouvementASupprimer.quantite}</p>
        <p class="delete-warning">Cette action ne modifiera pas le stock actuel.</p>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => mouvementASupprimer = null}>Annuler</button>
        <button class="btn-delete-movement" on:click={confirmerSuppressionMouvement}>Supprimer</button>
      </div>
    </div>
  </div>
{/if}

{#if modalEntree}
  <div class="overlay" on:click={() => modalEntree = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>+ Nouvelle Entrée de Stock</h2>
        <button class="close-btn" on:click={() => modalEntree = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Sélectionner l'article :</label>
          <select bind:value={formEntree.produitId}>
            {#each produitsListe as p}
              <option value={p.id}>{p.designation} (Stock actuel : {p.stockActuel})</option>
            {/each}
          </select>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Quantité reçue :</label>
            <input type="number" min="1" bind:value={formEntree.quantite} required />
          </div>
          <div class="form-group">
            <label>N° Bon / Document :</label>
            <input type="text" placeholder="Ex: BL-4590" bind:value={formEntree.referenceDocument} />
          </div>
        </div>
        <div class="form-group">
          <label>Motif de l'entrée :</label>
          <input type="text" bind:value={formEntree.motif} placeholder="Ex: Livraison fournisseur" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalEntree = false}>Annuler</button>
        <button class="btn-green" on:click={validerEntree}>Valider l'entrée</button>
      </div>
    </div>
  </div>
{/if}

{#if modalSortie}
  <div class="overlay" on:click={() => modalSortie = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>- Nouvelle Sortie de Stock</h2>
        <button class="close-btn" on:click={() => modalSortie = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Sélectionner l'article :</label>
          <select bind:value={formSortie.produitId}>
            {#each produitsListe as p}
              <option value={p.id}>{p.designation} (Stock actuel : {p.stockActuel})</option>
            {/each}
          </select>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Quantité retirée :</label>
            <input type="number" min="1" bind:value={formSortie.quantite} required />
          </div>
          <div class="form-group">
            <label>Référence interne :</label>
            <input type="text" placeholder="Ex: RET-001" bind:value={formSortie.referenceDocument} />
          </div>
        </div>
        <div class="form-group">
          <label>Motif de la sortie :</label>
          <input type="text" bind:value={formSortie.motif} placeholder="Ex: Casse matériel, perte, usage interne" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalSortie = false}>Annuler</button>
        <button class="btn-red" on:click={validerSortie}>Valider la sortie</button>
      </div>
    </div>
  </div>
{/if}

{#if modalAjustement}
  <div class="overlay" on:click={() => modalAjustement = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Ajustement d'Inventaire</h2>
        <button class="close-btn" on:click={() => modalAjustement = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>Sélectionner l'article :</label>
          <select bind:value={formAjustement.produitId} on:change={onChangementProduitAjustement}>
            {#each produitsListe as p}
              <option value={p.id}>{p.designation} (En base : {p.stockActuel})</option>
            {/each}
          </select>
        </div>
        <div class="form-group">
          <label>Nouveau stock réel compté :</label>
          <input type="number" bind:value={formAjustement.nouveauStock} required />
        </div>
        <div class="form-group">
          <label>Motif / Justification :</label>
          <input type="text" bind:value={formAjustement.motif} placeholder="Ex: Inventaire annuel, régularisation" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalAjustement = false}>Annuler</button>
        <button class="btn-primary" on:click={validerAjustement}>Appliquer l'ajustement</button>
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
    font-size: 0.9rem;
    transition: opacity 0.15s;
  }
  .btn:hover {
    opacity: 0.9;
  }
  .btn-entree {
    background: #16a34a;
    color: white;
  }
  .btn-sortie {
    background: #dc2626;
    color: white;
  }
  .btn-ajust {
    background: #475569;
    color: white;
  }
  .alert-box {
    background: #fffbeb;
    border: 1px solid #fef3c7;
    border-radius: 10px;
    padding: 16px;
    margin-bottom: 20px;
  }
  .alert-header {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #92400e;
    margin-bottom: 12px;
  }
  .alert-list {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  .alert-pill {
    background: #ffffff;
    border: 1px solid #fde68a;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 0.85rem;
    display: flex;
    gap: 8px;
  }
  .alert-pill .name {
    font-weight: 600;
    color: #78350f;
  }
  .alert-pill .stock-info {
    color: #b45309;
  }
  .card {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
    overflow: hidden;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #f1f5f9;
  }
  .card-header h2 {
    margin: 0;
    font-size: 1.15rem;
    color: #1e293b;
  }
  .btn-refresh {
    background: #f1f5f9;
    border: 1px solid #cbd5e1;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 500;
  }
  .btn-delete-movement {
    background: #fee2e2;
    border: 1px solid #fecaca;
    color: #991b1b;
    padding: 6px 10px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
    white-space: nowrap;
  }
  .btn-delete-movement:hover {
    background: #fecaca;
  }
  .delete-warning {
    color: #b45309;
    font-weight: 600;
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
  .text-center {
    text-align: center;
  }
  .font-bold {
    font-weight: 700;
  }
  .text-green {
    color: #16a34a;
  }
  .text-red {
    color: #dc2626;
  }
  .text-muted {
    color: #94a3b8;
  }
  .doc-ref {
    color: #64748b;
    font-size: 0.85rem;
    margin-left: 4px;
  }
  .badge {
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: 600;
  }
  .badge-entree {
    background: #dcfce7;
    color: #166534;
  }
  .badge-sortie {
    background: #fee2e2;
    color: #991b1b;
  }
  .badge-vente {
    background: #e0f2fe;
    color: #0369a1;
  }
  .badge-ajust {
    background: #f3e8ff;
    color: #6b21a8;
  }
  .empty {
    text-align: center;
    color: #94a3b8;
    padding: 36px !important;
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
    width: 500px;
    max-width: 90vw;
    box-shadow: 0 20px 25px -5px rgba(0,0,0,0.2);
    overflow: hidden;
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
  .btn-primary {
    background: #0284c7;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-green {
    background: #16a34a;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-red {
    background: #dc2626;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-secondary {
    background: #f1f5f9;
    color: #475569;
    border: 1px solid #cbd5e1;
    padding: 10px 18px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
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
