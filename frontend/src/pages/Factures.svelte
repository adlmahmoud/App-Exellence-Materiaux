<script>
  import { onMount } from 'svelte';
  import {
    ListerFactures,
    ImprimerFacturePDF,
    CreerFacture,
    ListerBons,
    ImprimerBonPDF,
    CreerBon,
    ListerProduits
  } from '../../wailsjs/go/main/App.js';

  let ongletActif = 'factures';
  let factures = [];
  let bons = [];
  let totalFactures = 0;
  let totalBons = 0;
  let page = 1;
  let limit = 20;
  let recherche = '';
  let chargement = false;
  let messageSucces = '';
  let messageErreur = '';

  let produitsDisponibles = [];

  let modalNouvelleFacture = false;
  let modalNouveauBon = false;

  let formFacture = {
    clientNom: '',
    clientAdresse: '',
    clientSirenSiret: '',
    modeReglement: 'especes',
    dateFacture: '',
    dateEcheance: '',
    conditionsReglement: 'Paiement comptant à réception de facture.',
    notes: '',
    lignes: []
  };

  let formBon = {
    typeBon: 'livraison',
    clientNom: '',
    referenceExterne: '',
    dateBon: '',
    notes: '',
    lignes: []
  };

  onMount(() => {
    chargerDonnees();
    chargerCatalogue();
  });

  async function chargerDonnees() {
    chargement = true;
    messageErreur = '';
    try {
      if (ongletActif === 'factures') {
        const res = await ListerFactures({
          recherche: recherche.trim(),
          page: page,
          limit: limit
        });
        factures = res.factures || [];
        totalFactures = res.total || 0;
      } else {
        const res = await ListerBons({
          recherche: recherche.trim(),
          page: page,
          limit: limit
        });
        bons = res.bons || [];
        totalBons = res.total || 0;
      }
    } catch (err) {
      messageErreur = 'Erreur lors du chargement : ' + err;
    } finally {
      chargement = false;
    }
  }

  async function chargerCatalogue() {
    try {
      const res = await ListerProduits({ limit: 100, page: 1, triChamp: 'designation', triAscendant: true });
      produitsDisponibles = res.produits || [];
    } catch (e) {}
  }

  let debounceTimer;

  function rechercheReactive() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      page = 1;
      chargerDonnees();
    }, 150);
  }

  function effacerRecherche() {
    recherche = '';
    page = 1;
    chargerDonnees();
  }

  function changerOnglet(onglet) {
    ongletActif = onglet;
    recherche = '';
    page = 1;
    chargerDonnees();
  }

  function ouvrirModalNouvelleFacture() {
    const today = new Date().toISOString().split('T')[0];
    const frToday = today.split('-').reverse().join('/');
    formFacture = {
      clientNom: '',
      clientAdresse: '',
      clientSirenSiret: '',
      modeReglement: 'especes',
      dateFacture: frToday,
      dateEcheance: frToday,
      conditionsReglement: 'Paiement comptant à réception de facture.',
      notes: '',
      lignes: []
    };
    ajouterLigneFacture();
    messageSucces = '';
    messageErreur = '';
    modalNouvelleFacture = true;
  }

  function ajouterLigneFacture() {
    formFacture.lignes = [
      ...formFacture.lignes,
      {
        produitId: '',
        designation: '',
        reference: '',
        codeBarres: '',
        quantite: 1,
        prixUnitaireHT: 0,
        tauxTva: 2000
      }
    ];
  }

  function onSelectionProduitFacture(index, e) {
    const prodId = e.target.value;
    const p = produitsDisponibles.find(item => item.id === prodId);
    if (p) {
      formFacture.lignes[index].produitId = p.id;
      formFacture.lignes[index].designation = p.designation;
      formFacture.lignes[index].reference = p.reference;
      formFacture.lignes[index].codeBarres = p.codeBarres;
      formFacture.lignes[index].prixUnitaireHT = (p.prixVenteHt / 100).toFixed(2);
      formFacture.lignes[index].tauxTva = p.tauxTva;
      formFacture.lignes = [...formFacture.lignes];
    }
  }

  function supprimerLigneFacture(index) {
    formFacture.lignes.splice(index, 1);
    if (formFacture.lignes.length === 0) {
      ajouterLigneFacture();
    } else {
      formFacture.lignes = [...formFacture.lignes];
    }
  }

  $: totalFactureHT = formFacture.lignes.reduce((acc, l) => acc + (Number(l.prixUnitaireHT || 0) * Number(l.quantite || 0)), 0);
  $: totalFactureTTC = formFacture.lignes.reduce((acc, l) => {
    const ht = Number(l.prixUnitaireHT || 0) * Number(l.quantite || 0);
    const tva = (ht * Number(l.tauxTva || 2000)) / 10000;
    return acc + ht + tva;
  }, 0);

  async function validerNouvelleFacture() {
    messageErreur = '';
    if (!formFacture.clientNom.trim()) {
      messageErreur = 'Le nom du client est requis.';
      return;
    }
    if (formFacture.lignes.length === 0) {
      messageErreur = 'Veuillez ajouter au moins un article.';
      return;
    }

    try {
      const payloadLignes = formFacture.lignes.map(l => ({
        produitId: l.produitId || null,
        codeBarres: l.codeBarres || '',
        reference: l.reference || '',
        designation: l.designation.trim() || 'Article sans désignation',
        quantite: Math.round(Number(l.quantite)),
        prixUnitaireHt: Math.round(Number(l.prixUnitaireHT) * 100),
        tauxTva: Number(l.tauxTva)
      }));

      const res = await CreerFacture({
        clientNom: formFacture.clientNom.trim(),
        clientAdresse: formFacture.clientAdresse.trim(),
        clientSirenSiret: formFacture.clientSirenSiret.trim(),
        modeReglement: formFacture.modeReglement,
        dateFacture: formFacture.dateFacture,
        dateEcheance: formFacture.dateEcheance,
        conditionsReglement: formFacture.conditionsReglement,
        notes: formFacture.notes,
        lignes: payloadLignes
      });

      modalNouvelleFacture = false;
      messageSucces = `Facture ${res.facture.numero} créée avec succès !`;
      chargerDonnees();
      imprimerFacture(res.facture.id);
    } catch (err) {
      messageErreur = String(err);
    }
  }

  function ouvrirModalNouveauBon() {
    const today = new Date().toISOString().split('T')[0];
    const frToday = today.split('-').reverse().join('/');
    formBon = {
      typeBon: 'livraison',
      clientNom: '',
      referenceExterne: '',
      dateBon: frToday,
      notes: '',
      lignes: []
    };
    ajouterLigneBon();
    messageSucces = '';
    messageErreur = '';
    modalNouveauBon = true;
  }

  function ajouterLigneBon() {
    formBon.lignes = [
      ...formBon.lignes,
      {
        produitId: '',
        designation: '',
        reference: '',
        codeBarres: '',
        quantite: 1,
        prixUnitaireHT: 0,
        tauxTva: 2000
      }
    ];
  }

  function onSelectionProduitBon(index, e) {
    const prodId = e.target.value;
    const p = produitsDisponibles.find(item => item.id === prodId);
    if (p) {
      formBon.lignes[index].produitId = p.id;
      formBon.lignes[index].designation = p.designation;
      formBon.lignes[index].reference = p.reference;
      formBon.lignes[index].codeBarres = p.codeBarres;
      formBon.lignes[index].prixUnitaireHT = (p.prixVenteHt / 100).toFixed(2);
      formBon.lignes[index].tauxTva = p.tauxTva;
      formBon.lignes = [...formBon.lignes];
    }
  }

  function supprimerLigneBon(index) {
    formBon.lignes.splice(index, 1);
    if (formBon.lignes.length === 0) {
      ajouterLigneBon();
    } else {
      formBon.lignes = [...formBon.lignes];
    }
  }

  async function validerNouveauBon() {
    messageErreur = '';
    if (!formBon.clientNom.trim()) {
      messageErreur = 'Le nom du destinataire est requis.';
      return;
    }

    try {
      const payloadLignes = formBon.lignes.map(l => ({
        produitId: l.produitId || null,
        codeBarres: l.codeBarres || '',
        reference: l.reference || '',
        designation: l.designation.trim() || 'Article',
        quantite: Math.round(Number(l.quantite)),
        prixUnitaireHt: Math.round(Number(l.prixUnitaireHT) * 100),
        tauxTva: Number(l.tauxTva)
      }));

      const res = await CreerBon({
        typeBon: formBon.typeBon,
        clientNom: formBon.clientNom.trim(),
        referenceExterne: formBon.referenceExterne.trim(),
        dateBon: formBon.dateBon,
        notes: formBon.notes,
        lignes: payloadLignes
      });

      modalNouveauBon = false;
      messageSucces = `Bon ${res.bon.numero} créé avec succès !`;
      chargerDonnees();
      imprimerBon(res.bon.id);
    } catch (err) {
      messageErreur = String(err);
    }
  }

  async function imprimerFacture(id) {
    try {
      await ImprimerFacturePDF(id);
    } catch (err) {
      messageErreur = 'Erreur impression : ' + err;
    }
  }

  async function imprimerBon(id) {
    try {
      await ImprimerBonPDF(id);
    } catch (err) {
      messageErreur = 'Erreur impression : ' + err;
    }
  }
</script>

<div class="page-container">
  <div class="header">
    <div>
      <h1>Factures & Bons de livraison</h1>
      <p class="desc">Numérotation séquentielle certifiée sans trou et impression PDF conforme</p>
    </div>
    <div class="actions">
      <button class="btn btn-primary" on:click={ouvrirModalNouvelleFacture}>+ Nouvelle Facture</button>
      <button class="btn btn-secondary" on:click={ouvrirModalNouveauBon}>+ Nouveau Bon</button>
    </div>
  </div>

  {#if messageSucces}
    <div class="alert success">{messageSucces}</div>
  {/if}
  {#if messageErreur}
    <div class="alert error">{messageErreur}</div>
  {/if}

  <div class="tabs-bar">
    <button
      class="tab-btn {ongletActif === 'factures' ? 'active' : ''}"
      on:click={() => changerOnglet('factures')}
    >
      Factures ({totalFactures})
    </button>
    <button
      class="tab-btn {ongletActif === 'bons' ? 'active' : ''}"
      on:click={() => changerOnglet('bons')}
    >
      Bons de livraison ({totalBons})
    </button>
  </div>

  <div class="search-bar">
    <div class="search-input-wrapper">
      <span class="search-icon">🔍</span>
      <input
        type="text"
        placeholder={ongletActif === 'factures' ? 'Rechercher une facture par numéro ou client...' : 'Rechercher un bon par numéro, destinataire ou réf...'}
        bind:value={recherche}
        on:input={rechercheReactive}
      />
      {#if recherche}
        <button class="btn-clear" on:click={effacerRecherche} title="Effacer">&times;</button>
      {/if}
    </div>
  </div>

  <div class="card">
    {#if ongletActif === 'factures'}
      <table class="table">
        <thead>
          <tr>
            <th>Numéro</th>
            <th>Date</th>
            <th>Client</th>
            <th>Règlement</th>
            <th class="text-right">Total HT</th>
            <th class="text-right">Total TTC</th>
            <th>Statut</th>
            <th class="text-right">Action</th>
          </tr>
        </thead>
        <tbody>
          {#if chargement}
            <tr><td colspan="8" class="empty">Chargement des factures...</td></tr>
          {:else if factures.length === 0}
            <tr><td colspan="8" class="empty">Aucune facture enregistrée pour le moment.</td></tr>
          {:else}
            {#each factures as f}
              <tr>
                <td><strong>{f.numero}</strong></td>
                <td>{f.dateFacture}</td>
                <td>{f.clientNom}</td>
                <td>
                  <span class="badge-mode">{f.modeReglement.toUpperCase()}</span>
                </td>
                <td class="text-right">{(f.totalHt / 100).toFixed(2)} €</td>
                <td class="text-right font-bold text-blue">{(f.totalTtc / 100).toFixed(2)} €</td>
                <td>
                  <span class="badge badge-success">Validée</span>
                </td>
                <td class="text-right">
                  <button class="btn-pdf" on:click={() => imprimerFacture(f.id)}>
                    📄 Imprimer PDF
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    {:else}
      <table class="table">
        <thead>
          <tr>
            <th>Numéro</th>
            <th>Type</th>
            <th>Date</th>
            <th>Destinataire</th>
            <th>Réf. Externe</th>
            <th class="text-right">Articles</th>
            <th>Statut</th>
            <th class="text-right">Action</th>
          </tr>
        </thead>
        <tbody>
          {#if chargement}
            <tr><td colspan="8" class="empty">Chargement des bons...</td></tr>
          {:else if bons.length === 0}
            <tr><td colspan="8" class="empty">Aucun bon enregistré pour le moment.</td></tr>
          {:else}
            {#each bons as b}
              <tr>
                <td><strong>{b.numero}</strong></td>
                <td>
                  <span class="badge badge-bon">{b.typeBon}</span>
                </td>
                <td>{b.dateBon}</td>
                <td>{b.clientNom}</td>
                <td>{b.referenceExterne || '-'}</td>
                <td class="text-right font-bold">{(b.totalTtc / 100).toFixed(2)} €</td>
                <td>
                  <span class="badge badge-info">{b.statut}</span>
                </td>
                <td class="text-right">
                  <button class="btn-pdf" on:click={() => imprimerBon(b.id)}>
                    📄 Imprimer PDF
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    {/if}
  </div>
</div>

{#if modalNouvelleFacture}
  <div class="overlay" on:click={() => modalNouvelleFacture = false}>
    <div class="modal modal-large" on:click|stopPropagation>
      <div class="modal-header">
        <h2>+ Nouvelle Facture</h2>
        <button class="close-btn" on:click={() => modalNouvelleFacture = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-row">
          <div class="form-group flex-2">
            <label>Nom ou Société du client * :</label>
            <input type="text" placeholder="Ex: Entreprise Martin Bâtiment" bind:value={formFacture.clientNom} required />
          </div>
          <div class="form-group flex-1">
            <label>SIRET (si professionnel) :</label>
            <input type="text" placeholder="14 chiffres" bind:value={formFacture.clientSirenSiret} />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group flex-2">
            <label>Adresse complète :</label>
            <input type="text" placeholder="Ex: 12 Rue du Progrès, 95500 Gonesse" bind:value={formFacture.clientAdresse} />
          </div>
          <div class="form-group flex-1">
            <label>Mode de règlement :</label>
            <select bind:value={formFacture.modeReglement}>
              <option value="especes">Espèces</option>
              <option value="cb">Carte Bancaire</option>
              <option value="virement">Virement Bancaire</option>
              <option value="cheque">Chèque</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>Date de facture :</label>
            <input type="text" bind:value={formFacture.dateFacture} />
          </div>
          <div class="form-group">
            <label>Date d'échéance :</label>
            <input type="text" bind:value={formFacture.dateEcheance} />
          </div>
        </div>

        <div class="lignes-section">
          <div class="lignes-header">
            <strong>Articles facturés</strong>
            <button class="btn-ajouter-ligne" on:click={ajouterLigneFacture}>+ Ajouter une ligne</button>
          </div>

          <table class="table-lignes">
            <thead>
              <tr>
                <th style="width: 25%;">Article rapide</th>
                <th>Désignation</th>
                <th style="width: 10%;">Qté</th>
                <th style="width: 15%;">PU HT (€)</th>
                <th style="width: 12%;">TVA</th>
                <th style="width: 15%;">Total HT</th>
                <th style="width: 5%;"></th>
              </tr>
            </thead>
            <tbody>
              {#each formFacture.lignes as l, idx}
                <tr>
                  <td>
                    <select on:change={(e) => onSelectionProduitFacture(idx, e)}>
                      <option value="">-- Choisir --</option>
                      {#each produitsDisponibles as p}
                        <option value={p.id}>{p.designation}</option>
                      {/each}
                    </select>
                  </td>
                  <td>
                    <input type="text" bind:value={l.designation} placeholder="Désignation de l'article" required />
                  </td>
                  <td>
                    <input type="number" min="1" bind:value={l.quantite} />
                  </td>
                  <td>
                    <input type="number" step="0.01" bind:value={l.prixUnitaireHT} />
                  </td>
                  <td>
                    <select bind:value={l.tauxTva}>
                      <option value={2000}>20 %</option>
                      <option value={1000}>10 %</option>
                      <option value={550}>5.5 %</option>
                      <option value={0}>0 %</option>
                    </select>
                  </td>
                  <td>
                    <strong>{(Number(l.prixUnitaireHT || 0) * Number(l.quantite || 0)).toFixed(2)} €</strong>
                  </td>
                  <td>
                    <button class="btn-del" on:click={() => supprimerLigneFacture(idx)}>&times;</button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>

        <div class="recap-facture">
          <div class="total-line">
            <span>Total HT :</span>
            <strong>{totalFactureHT.toFixed(2)} €</strong>
          </div>
          <div class="total-line ttc">
            <span>TOTAL TTC :</span>
            <strong>{totalFactureTTC.toFixed(2)} €</strong>
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalNouvelleFacture = false}>Annuler</button>
        <button class="btn-primary" on:click={validerNouvelleFacture}>Générer la Facture & Imprimer PDF</button>
      </div>
    </div>
  </div>
{/if}

{#if modalNouveauBon}
  <div class="overlay" on:click={() => modalNouveauBon = false}>
    <div class="modal modal-large" on:click|stopPropagation>
      <div class="modal-header">
        <h2>+ Nouveau Bon de livraison</h2>
        <button class="close-btn" on:click={() => modalNouveauBon = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-row">
          <div class="form-group flex-2">
            <label>Client destinataire * :</label>
            <input type="text" placeholder="Ex: Chantier BTP Nord" bind:value={formBon.clientNom} required />
          </div>
          <div class="form-group flex-1">
            <label>Type de document :</label>
            <select bind:value={formBon.typeBon}>
              <option value="livraison">Bon de livraison</option>
              <option value="commande">Bon de commande</option>
              <option value="devis">Devis</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group flex-2">
            <label>Référence externe / N° Commande :</label>
            <input type="text" placeholder="Ex: CMD-8812" bind:value={formBon.referenceExterne} />
          </div>
          <div class="form-group flex-1">
            <label>Date du document :</label>
            <input type="text" bind:value={formBon.dateBon} />
          </div>
        </div>

        <div class="lignes-section">
          <div class="lignes-header">
            <strong>Articles à livrer</strong>
            <button class="btn-ajouter-ligne" on:click={ajouterLigneBon}>+ Ajouter une ligne</button>
          </div>

          <table class="table-lignes">
            <thead>
              <tr>
                <th style="width: 30%;">Article rapide</th>
                <th>Désignation</th>
                <th style="width: 15%;">Quantité</th>
                <th style="width: 5%;"></th>
              </tr>
            </thead>
            <tbody>
              {#each formBon.lignes as l, idx}
                <tr>
                  <td>
                    <select on:change={(e) => onSelectionProduitBon(idx, e)}>
                      <option value="">-- Choisir --</option>
                      {#each produitsDisponibles as p}
                        <option value={p.id}>{p.designation}</option>
                      {/each}
                    </select>
                  </td>
                  <td>
                    <input type="text" bind:value={l.designation} placeholder="Désignation" required />
                  </td>
                  <td>
                    <input type="number" min="1" bind:value={l.quantite} />
                  </td>
                  <td>
                    <button class="btn-del" on:click={() => supprimerLigneBon(idx)}>&times;</button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" on:click={() => modalNouveauBon = false}>Annuler</button>
        <button class="btn-primary" on:click={validerNouveauBon}>Générer le Bon & Imprimer PDF</button>
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
  }
  .btn-primary {
    background: #0284c7;
    color: white;
  }
  .btn-secondary {
    background: #475569;
    color: white;
  }
  .tabs-bar {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
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
    font-size: 0.95rem;
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
  .tab-btn {
    padding: 10px 18px;
    background: #e2e8f0;
    color: #475569;
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .tab-btn.active {
    background: #0284c7;
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
  .font-bold {
    font-weight: 700;
  }
  .text-blue {
    color: #0284c7;
  }
  .empty {
    text-align: center;
    color: #94a3b8;
    padding: 40px !important;
  }
  .badge {
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: 600;
  }
  .badge-success {
    background: #dcfce7;
    color: #166534;
  }
  .badge-bon {
    background: #e0e7ff;
    color: #3730a3;
  }
  .badge-info {
    background: #f1f5f9;
    color: #475569;
  }
  .badge-mode {
    background: #f8fafc;
    border: 1px solid #cbd5e1;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 600;
    color: #475569;
  }
  .btn-pdf {
    background: #f1f5f9;
    border: 1px solid #cbd5e1;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
    color: #1e293b;
    transition: background 0.15s;
  }
  .btn-pdf:hover {
    background: #e2e8f0;
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
  .modal-large {
    width: 800px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
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
    overflow-y: auto;
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
    margin-bottom: 12px;
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
  .lignes-section {
    margin-top: 16px;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 12px;
    background: #fafafa;
  }
  .lignes-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
  }
  .btn-ajouter-ligne {
    background: #334155;
    color: white;
    border: none;
    padding: 6px 12px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
  }
  .table-lignes {
    width: 100%;
    border-collapse: collapse;
  }
  .table-lignes th, .table-lignes td {
    padding: 6px 8px;
    font-size: 0.85rem;
  }
  .table-lignes th {
    text-align: left;
    color: #64748b;
  }
  .btn-del {
    background: transparent;
    border: none;
    color: #dc2626;
    font-size: 1.2rem;
    cursor: pointer;
    font-weight: bold;
  }
  .recap-facture {
    margin-top: 16px;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
  }
  .total-line {
    display: flex;
    gap: 20px;
    font-size: 0.95rem;
    color: #475569;
  }
  .total-line.ttc {
    font-size: 1.15rem;
    color: #0f172a;
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
