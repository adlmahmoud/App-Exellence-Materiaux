<script>
  import { API } from '../lib/api.js';
  import { onMount } from 'svelte';
  const {
    ChangeAccessPasswords,
    GetParametres,
    UpdateParametres,
    ImporterLogo,
    SauvegarderBase,
    GetConfigTPE,
    SauvegarderConfigTPE,
    TesterConnexionTPE,
    ListerPortsCOMTPE
  } = API;

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
  let capitalEuros = '';
  let enregistrementEnCours = false;
  let currentAdminPassword = '';
  let newAdminPassword = '';
  let confirmAdminPassword = '';
  let newUserPassword = '';
  let confirmUserPassword = '';

  let configTPE = {
    mode: 'simulateur',
    portCom: 'COM1',
    baudRate: 9600,
    adresseIp: '192.168.1.50',
    portReseau: 8888,
    timeoutSec: 45,
    simulateurScenario: 'succes',
    simulateurDelaiSec: 2
  };
  let portsCOM = [];
  let testTPEResultat = null;
  let testTPEEnCours = false;
  let tpeMessage = '';
  let tpeMessageErreur = '';

  onMount(async () => {
    try {
      const res = await GetParametres();
      if (res) {
        params = { ...params, ...res };
        capitalEuros = (res.capital / 100).toFixed(2);
      }
      const cfg = await GetConfigTPE();
      if (cfg && cfg.mode) {
        configTPE = { ...configTPE, ...cfg };
      }
      portsCOM = (await ListerPortsCOMTPE()) || [];
    } catch (err) {
      messageErreur = 'Erreur lors du chargement des paramètres.';
    } finally {
      chargement = false;
    }
  });

  async function enregistrerInformations() {
    message = '';
    messageErreur = '';
    enregistrementEnCours = true;
    try {
      const valNettoyee = String(capitalEuros ?? '').trim().replace(/\s+/g, '').replace(',', '.');
      const capFloat = parseFloat(valNettoyee) || 0;
      const capitalCentimes = Math.round(capFloat * 100);
      const payload = {
        denomination: params.denomination || '',
        formeJuridique: params.formeJuridique || '',
        capital: capitalCentimes,
        adresse: params.adresse || '',
        codePostal: params.codePostal || '',
        ville: params.ville || '',
        siren: params.siren || '',
        siret: params.siret || '',
        rcs: params.rcs || '',
        numeroTva: params.numeroTva || '',
        telephone: params.telephone || '',
        email: params.email || '',
        siteWeb: params.siteWeb || '',
        logoPath: params.logoPath || '',
        devise: params.devise || 'EUR',
        tauxTvaDefaut: Number(params.tauxTvaDefaut) || 2000,
        mentionsLegales: params.mentionsLegales || '',
        conditionsReglement: params.conditionsReglement || ''
      };
      await UpdateParametres(payload);
      params.capital = capitalCentimes;
      capitalEuros = (capitalCentimes / 100).toFixed(2);
      message = 'Les informations de l’entreprise ont été enregistrées avec succès.';
    } catch (err) {
      messageErreur = 'Erreur lors de l’enregistrement : ' + (err?.message || err);
    } finally {
      enregistrementEnCours = false;
    }
  }

  async function testerTPE() {
    testTPEEnCours = true;
    testTPEResultat = null;
    tpeMessage = '';
    tpeMessageErreur = '';
    try {
      await SauvegarderConfigTPE(configTPE);
      const res = await TesterConnexionTPE();
      testTPEResultat = res;
    } catch (e) {
      tpeMessageErreur = 'Erreur test TPE: ' + e;
    } finally {
      testTPEEnCours = false;
    }
  }

  async function enregistrerTPE() {
    tpeMessage = '';
    tpeMessageErreur = '';
    try {
      await SauvegarderConfigTPE(configTPE);
      tpeMessage = 'Configuration du TPE enregistrée avec succès !';
    } catch (e) {
      tpeMessageErreur = 'Erreur sauvegarde TPE: ' + e;
    }
  }

  async function actualiserPortsCOM() {
    try {
      portsCOM = (await ListerPortsCOMTPE()) || [];
    } catch (e) {}
  }

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

  async function changerMotsDePasse() {
    message = '';
    messageErreur = '';
    if (newAdminPassword.length < 8 || newUserPassword.length < 8) {
      messageErreur = 'Chaque nouveau mot de passe doit contenir au moins 8 caractères.';
      return;
    }
    if (newAdminPassword === newUserPassword) {
      messageErreur = 'Les mots de passe administrateur et utilisateur doivent être différents.';
      return;
    }
    if (newAdminPassword !== confirmAdminPassword || newUserPassword !== confirmUserPassword) {
      messageErreur = 'Les confirmations des mots de passe ne correspondent pas.';
      return;
    }

    try {
      await ChangeAccessPasswords(currentAdminPassword, newAdminPassword, newUserPassword);
      currentAdminPassword = '';
      newAdminPassword = '';
      confirmAdminPassword = '';
      newUserPassword = '';
      confirmUserPassword = '';
      message = 'Les mots de passe ont été modifiés.';
    } catch (err) {
      messageErreur = String(err);
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

    <form class="card access-card" on:submit|preventDefault={changerMotsDePasse}>
      <h2>Accès administrateur et utilisateur</h2>
      <div class="form-row">
        <div class="form-group">
          <label for="current-admin-password">Mot de passe administrateur actuel :</label>
          <input id="current-admin-password" type="password" bind:value={currentAdminPassword} autocomplete="current-password" required />
        </div>
        <div class="form-group">
          <label for="new-admin-password">Nouveau mot de passe administrateur :</label>
          <input id="new-admin-password" type="password" minlength="8" bind:value={newAdminPassword} autocomplete="new-password" required />
        </div>
        <div class="form-group">
          <label for="confirm-admin-password">Confirmer le mot de passe administrateur :</label>
          <input id="confirm-admin-password" type="password" minlength="8" bind:value={confirmAdminPassword} autocomplete="new-password" required />
        </div>
        <div class="form-group">
          <label for="new-user-password">Nouveau mot de passe utilisateur :</label>
          <input id="new-user-password" type="password" minlength="8" bind:value={newUserPassword} autocomplete="new-password" required />
        </div>
        <div class="form-group">
          <label for="confirm-user-password">Confirmer le mot de passe utilisateur :</label>
          <input id="confirm-user-password" type="password" minlength="8" bind:value={confirmUserPassword} autocomplete="new-password" required />
        </div>
      </div>
      <button class="btn-backup" type="submit">Modifier les mots de passe</button>
    </form>

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
            <input type="text" inputmode="decimal" bind:value={capitalEuros} placeholder="ex: 1000.00" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label>Téléphone :</label>
            <input type="text" bind:value={params.telephone} />
          </div>
          <div class="form-group">
            <label>Email :</label>
            <input type="text" bind:value={params.email} />
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
        <div style="margin-top: 14px; margin-bottom: 20px;">
          <button type="button" class="btn-save-tpe" style="width: 100%;" on:click={enregistrerInformations} disabled={enregistrementEnCours}>
            {enregistrementEnCours ? 'Enregistrement en cours...' : '💾 Enregistrer les informations de l’entreprise'}
          </button>
        </div>
        <div class="backup-box">
          <div>
            <strong>Sauvegarde locale SQLite</strong>
            <p>Crée une copie instantanée et sécurisée de la base de données.</p>
          </div>
          <button class="btn-backup" on:click={faireSauvegarde}>Sauvegarder la base</button>
        </div>
      </div>

      <div class="card full-width">
        <h2>💳 Terminal de Paiement Électronique (TPE)</h2>
        <p class="section-desc">Configuration universelle du terminal de paiement par carte bancaire branché en physique ou en réseau.</p>

        {#if tpeMessage}
          <div class="alert success">{tpeMessage}</div>
        {/if}
        {#if tpeMessageErreur}
          <div class="alert error">{tpeMessageErreur}</div>
        {/if}

        <div class="form-group">
          <label>Mode de connexion au TPE :</label>
          <div class="tpe-mode-selector">
            <button
              type="button"
              class="btn-mode {configTPE.mode === 'simulateur' ? 'active' : ''}"
              on:click={() => configTPE.mode = 'simulateur'}
            >
              🧪 Simulateur / Test (Sans matériel physique)
            </button>
            <button
              type="button"
              class="btn-mode {configTPE.mode === 'serie' ? 'active' : ''}"
              on:click={() => configTPE.mode = 'serie'}
            >
              🔌 Câble USB / Série (Port COM)
            </button>
            <button
              type="button"
              class="btn-mode {configTPE.mode === 'reseau' ? 'active' : ''}"
              on:click={() => configTPE.mode = 'reseau'}
            >
              🌐 Réseau Local (Ethernet / IP)
            </button>
          </div>
        </div>

        {#if configTPE.mode === 'simulateur'}
          <div class="tpe-settings-panel">
            <div class="form-row">
              <div class="form-group">
                <label>Résultat simulé lors des encaissements en caisse :</label>
                <select bind:value={configTPE.simulateurScenario}>
                  <option value="succes">✅ Paiement Validé (Succès immédiat)</option>
                  <option value="refus">❌ Paiement Refusé par la banque</option>
                  <option value="annulation">⚠️ Transaction Annulée sur le TPE</option>
                  <option value="aleatoire">🎲 Aléatoire (Réaliste)</option>
                </select>
              </div>
              <div class="form-group">
                <label>Délai de simulation (secondes) :</label>
                <input type="number" min="1" max="10" bind:value={configTPE.simulateurDelaiSec} />
              </div>
            </div>
            <p class="tpe-info-box">
              ℹ️ Ce mode vous permet de tester immédiatement tous les scénarios (validation, refus, annulation) sans aucun TPE branché.
            </p>
          </div>
        {:else if configTPE.mode === 'serie'}
          <div class="tpe-settings-panel">
            <div class="form-row">
              <div class="form-group">
                <label>Port COM :</label>
                <div class="com-row">
                  {#if portsCOM.length > 0}
                    <select bind:value={configTPE.portCom}>
                      {#each portsCOM as p}
                        <option value={p}>{p}</option>
                      {/each}
                    </select>
                  {:else}
                    <input type="text" bind:value={configTPE.portCom} placeholder="ex: COM3" />
                  {/if}
                  <button type="button" class="btn-refresh" on:click={actualiserPortsCOM} title="Actualiser les ports">
                    🔄
                  </button>
                </div>
              </div>
              <div class="form-group">
                <label>Vitesse (BaudRate) :</label>
                <select bind:value={configTPE.baudRate}>
                  <option value={9600}>9600 bauds (Standard Concert)</option>
                  <option value={19200}>19200 bauds</option>
                  <option value={38400}>38400 bauds</option>
                  <option value={115200}>115200 bauds (PAX / USB rapide)</option>
                </select>
              </div>
            </div>
            <p class="tpe-info-box">
              🔌 Protocole Concert V3 / CB PR V3.0 (Ingenico Desk/Move, PAX, Verifone reliés en USB).
            </p>
          </div>
        {:else}
          <div class="tpe-settings-panel">
            <div class="form-row">
              <div class="form-group">
                <label>Adresse IP du TPE :</label>
                <input type="text" bind:value={configTPE.adresseIp} placeholder="ex: 192.168.1.50" />
              </div>
              <div class="form-group">
                <label>Port réseau :</label>
                <input type="number" bind:value={configTPE.portReseau} placeholder="8888" />
              </div>
            </div>
            <p class="tpe-info-box">
              🌐 Protocole Concert sur socket TCP (TPE relié en câble Ethernet à votre box ou switch).
            </p>
          </div>
        {/if}

        <div class="form-row">
          <div class="form-group">
            <label>Délai d'attente maximum (secondes) :</label>
            <input type="number" min="15" max="120" bind:value={configTPE.timeoutSec} />
          </div>
        </div>

        {#if testTPEResultat}
          <div class="alert {testTPEResultat.succes ? 'success' : 'error'}">
            <strong>{testTPEResultat.succes ? '✅ Connexion réussie :' : '❌ Échec de connexion :'}</strong>
            <span>{testTPEResultat.message}</span>
            {#if testTPEResultat.details}
              <div style="font-size:0.85rem; margin-top:4px;">{testTPEResultat.details}</div>
            {/if}
          </div>
        {/if}

        <div class="tpe-footer-buttons">
          <button type="button" class="btn-test-tpe" on:click={testerTPE} disabled={testTPEEnCours}>
            {testTPEEnCours ? 'Test en cours...' : '🔍 Tester la communication avec le TPE'}
          </button>
          <button type="button" class="btn-save-tpe" on:click={enregistrerTPE}>
            💾 Enregistrer la configuration TPE
          </button>
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
  .section-desc {
    color: #64748b;
    font-size: 0.9rem;
    margin-top: -8px;
    margin-bottom: 16px;
  }
  .tpe-mode-selector {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
    margin-bottom: 16px;
  }
  .btn-mode {
    background: #f8fafc;
    border: 2px solid #cbd5e1;
    border-radius: 8px;
    padding: 12px;
    font-size: 0.85rem;
    font-weight: 600;
    color: #334155;
    cursor: pointer;
    transition: all 0.15s;
    text-align: center;
  }
  .btn-mode.active {
    background: #f0fdf4;
    border-color: #16a34a;
    color: #166534;
  }
  .tpe-settings-panel {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
  }
  .tpe-info-box {
    margin: 8px 0 0;
    font-size: 0.85rem;
    color: #475569;
  }
  .com-row {
    display: flex;
    gap: 8px;
  }
  .btn-refresh {
    background: #e2e8f0;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    padding: 0 12px;
    cursor: pointer;
  }
  .tpe-footer-buttons {
    display: flex;
    gap: 12px;
    margin-top: 16px;
  }
  .btn-test-tpe {
    background: #0284c7;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 12px 18px;
    font-weight: 600;
    cursor: pointer;
    flex: 1;
  }
  .btn-test-tpe:hover:not(:disabled) {
    background: #0369a1;
  }
  .btn-test-tpe:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .btn-save-tpe {
    background: #16a34a;
    color: white;
    border: none;
    border-radius: 6px;
    padding: 12px 18px;
    font-weight: 600;
    cursor: pointer;
    flex: 1;
  }
  .btn-save-tpe:hover {
    background: #15803d;
  }
</style>
