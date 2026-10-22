<script>
  import { GetParametres, ImporterLogo, SauvegarderBase } from '../../wailsjs/go/main/App.js';
  import { onMount } from 'svelte';

  let params = {
    denomination: '',
    formeJuridique: '',
    capital: 0,
    adresse: '',
    codePostal: '',
    ville: '',
    siren: '',
    siret: '',
    rcs: '',
    numeroTva: '',
    telephone: '',
    email: '',
    siteWeb: '',
    logoPath: '',
    devise: 'EUR',
    tauxTvaDefaut: 2000,
    mentionsLegales: '',
    conditionsReglement: ''
  };

  let chargement = true;
  let message = '';
  let messageErreur = '';

  onMount(async () => {
    try {
      const res = await GetParametres();
      if (res) {
        params = { ...params, ...res };
      }
    } catch (err) {
      messageErreur = 'Erreur lors du chargement des paramètres.';
    } finally {
      chargement = false;
    }
  });

  async function faireSauvegarde() {
    message = '';
    messageErreur = '';
    try {
      const res = await SauvegarderBase('.');
      message = 'Sauvegarde créée avec succès : ' + res;
    } catch (err) {
      messageErreur = 'Erreur de sauvegarde : ' + err;
    }
  }

  async function choisirLogo() {
    message = '';
    messageErreur = '';
    try {
      const path = await ImporterLogo();
      if (path) {
        params = { ...params, logoPath: path };
        message = 'Logo enregistré pour les factures et bons de livraison.';
      }
    } catch (err) {
      messageErreur = 'Erreur lors de l’ajout du logo : ' + err;
    }
  }

  function nomFichierLogo(path) {
    return path ? path.split(/[\\/]/).pop() : '';
  }
</script>

<div class="page-container">
  <div class="header">
    <h1>Paramètres de l'entreprise</h1>
    <p class="desc">Coordonnées légales, TVA et sauvegarde de la base locale</p>
  </div>

  {#if chargement}
    <div class="card loading">Chargement des données...</div>
  {:else}
    {#if message}
      <div class="alert success">{message}</div>
    {/if}
    {#if messageErreur}
      <div class="alert error">{messageErreur}</div>
    {/if}

    <div class="grid">
      <div class="card">
        <h2>Identité</h2>
        <div class="form-group">
          <label>Dénomination :</label>
          <input type="text" bind:value={params.denomination} />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Forme juridique :</label>
            <input type="text" bind:value={params.formeJuridique} />
          </div>
          <div class="form-group">
            <label>Capital (€) :</label>
            <input type="number" value={(params.capital / 100).toFixed(2)} readonly />
          </div>
        </div>
        <div class="form-group">
          <label>Adresse :</label>
          <input type="text" bind:value={params.adresse} />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Code Postal :</label>
            <input type="text" bind:value={params.codePostal} />
          </div>
          <div class="form-group">
            <label>Ville :</label>
            <input type="text" bind:value={params.ville} />
          </div>
        </div>
      </div>

      <div class="card">
        <h2>Immatriculation & Fiscalité</h2>
        <div class="form-row">
          <div class="form-group">
            <label>SIREN :</label>
            <input type="text" bind:value={params.siren} />
          </div>
          <div class="form-group">
            <label>RCS :</label>
            <input type="text" bind:value={params.rcs} />
          </div>
        </div>
        <div class="form-group">
          <label>N° TVA intracommunautaire :</label>
          <input type="text" bind:value={params.numeroTva} />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Devise :</label>
            <input type="text" bind:value={params.devise} />
          </div>
          <div class="form-group">
            <label>Taux TVA par défaut :</label>
            <input type="text" value={(params.tauxTvaDefaut / 100) + ' %'} readonly />
          </div>
        </div>
      </div>

      <div class="card full-width">
        <h2>Logo des documents</h2>
        <div class="logo-box">
          <div>
            <strong>{params.logoPath ? 'Logo configuré' : 'Aucun logo configuré'}</strong>
            <p>{params.logoPath ? nomFichierLogo(params.logoPath) : 'Le logo apparaîtra sur les factures et bons de livraison.'}</p>
          </div>
          <button class="btn-logo" on:click={choisirLogo}>
            <span aria-hidden="true">📷</span>
            {params.logoPath ? 'Remplacer le logo' : 'Ajouter un logo'}
          </button>
        </div>
      </div>

      <div class="card full-width">
        <h2>Mentions Facture & Sauvegarde</h2>
        <div class="form-group">
          <label>Conditions de règlement / mentions complémentaires :</label>
          <input type="text" bind:value={params.conditionsReglement} />
        </div>
        <div class="backup-box">
          <div>
            <strong>Sauvegarde locale SQLite</strong>
            <p>Crée une copie instantanée et sécurisée de la base de données.</p>
          </div>
          <button class="btn-backup" on:click={faireSauvegarde}>Sauvegarder la base</button>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .page-container {
    padding: 24px;
    max-width: 1100px;
    margin: 0 auto;
  }
  .header {
    margin-bottom: 24px;
  }
  .header h1 {
    margin: 0 0 6px;
    font-size: 1.6rem;
    color: #0f172a;
  }
  .desc {
    margin: 0;
    color: #64748b;
  }
  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }
  .full-width {
    grid-column: 1 / -1;
  }
  .card {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
    padding: 20px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  }
  .card h2 {
    margin: 0 0 16px;
    font-size: 1.1rem;
    color: #1e293b;
    border-bottom: 1px solid #f1f5f9;
    padding-bottom: 8px;
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
  input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 0.95rem;
    box-sizing: border-box;
  }
  input[readonly] {
    background: #f8fafc;
    color: #64748b;
  }
  .backup-box {
    margin-top: 16px;
    padding: 16px;
    background: #f8fafc;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border: 1px dashed #cbd5e1;
  }
  .logo-box {
    padding: 16px;
    background: #f8fafc;
    border-radius: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    border: 1px dashed #cbd5e1;
  }
  .logo-box p {
    margin: 4px 0 0;
    font-size: 0.85rem;
    color: #64748b;
  }
  .btn-logo {
    background: #0284c7;
    color: white;
    border: none;
    padding: 10px 18px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    white-space: nowrap;
  }
  .backup-box p {
    margin: 4px 0 0;
    font-size: 0.85rem;
    color: #64748b;
  }
  .btn-backup {
    background: #0284c7;
    color: white;
    border: none;
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
