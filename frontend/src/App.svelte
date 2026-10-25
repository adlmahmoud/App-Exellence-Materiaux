<script>
  import { onMount } from 'svelte';
  import logo from './assets/images/gestion-stock.png';
  import { API } from './lib/api.js';
  import Caisse from './pages/Caisse.svelte';
  import Produits from './pages/Produits.svelte';
  import Stock from './pages/Stock.svelte';
  import Factures from './pages/Factures.svelte';
  import Parametres from './pages/Parametres.svelte';
  import APropos from './pages/APropos.svelte';

  let pageActive = 'caisse';
  let modalAPropos = false;
  let authLoading = true;
  let setupRequired = false;
  let currentRole = '';
  let selectedRole = 'admin';
  let password = '';
  let adminPassword = '';
  let userPassword = '';
  let adminPasswordConfirm = '';
  let userPasswordConfirm = '';
  let authError = '';
  let authBusy = false;

  let serverInfo = { running: false, port: 8765, ips: [], mode: 'standalone', serverUrl: '' };
  let modalReseau = false;
  let ipSaisie = '';
  let erreurReseau = '';
  let reseauBusy = false;

  onMount(async () => {
    try {
      await actualiserInfoServeur();
      setupRequired = !(await API.GetAuthStatus());
      if (!setupRequired) {
        currentRole = await API.GetCurrentRole();
      }
    } catch (err) {
      authError = String(err);
    } finally {
      authLoading = false;
    }
  });

  async function actualiserInfoServeur() {
    try {
      serverInfo = await API.GetServerInfo();
      if (serverInfo && serverInfo.serverUrl) {
        ipSaisie = serverInfo.serverUrl.replace(/^https?:\/\//, '').replace(/:8765$/, '');
      }
    } catch (e) {}
  }

  function changerPage(p) {
    pageActive = p;
  }

  async function configurerMotsDePasse() {
    authError = '';
    if (adminPassword.length < 8 || userPassword.length < 8) {
      authError = 'Chaque mot de passe doit contenir au moins 8 caractères.';
      return;
    }
    if (adminPassword === userPassword) {
      authError = 'Les deux mots de passe doivent être différents.';
      return;
    }
    if (adminPassword !== adminPasswordConfirm || userPassword !== userPasswordConfirm) {
      authError = 'Les confirmations des mots de passe ne correspondent pas.';
      return;
    }

    authBusy = true;
    try {
      await API.SetupAccessPasswords(adminPassword, userPassword);
      setupRequired = false;
      adminPassword = '';
      userPassword = '';
      adminPasswordConfirm = '';
      userPasswordConfirm = '';
      authError = '';
    } catch (err) {
      authError = String(err);
    } finally {
      authBusy = false;
    }
  }

  async function ouvrirSession() {
    authError = '';
    authBusy = true;
    try {
      await API.Login(selectedRole, password);
      currentRole = selectedRole;
      password = '';
      pageActive = 'caisse';
      await actualiserInfoServeur();
    } catch (err) {
      authError = String(err);
    } finally {
      authBusy = false;
    }
  }

  async function fermerSession() {
    try {
      await API.Logout();
      currentRole = '';
      pageActive = 'caisse';
      authError = '';
      await actualiserInfoServeur();
    } catch (err) {
      authError = String(err);
    }
  }

  async function enregistrerLiaisonClient() {
    erreurReseau = '';
    reseauBusy = true;
    try {
      await API.SetNetworkMode('client', ipSaisie);
      modalReseau = false;
      await actualiserInfoServeur();
      setupRequired = !(await API.GetAuthStatus());
    } catch (err) {
      erreurReseau = String(err);
    } finally {
      reseauBusy = false;
    }
  }

  async function repasserEnModeLocal() {
    erreurReseau = '';
    reseauBusy = true;
    try {
      await API.SetNetworkMode('standalone', '');
      modalReseau = false;
      await actualiserInfoServeur();
      setupRequired = !(await API.GetAuthStatus());
    } catch (err) {
      erreurReseau = String(err);
    } finally {
      reseauBusy = false;
    }
  }
</script>

{#if authLoading}
  <div class="auth-screen"><div class="auth-card">Chargement...</div></div>
{:else if setupRequired}
  <div class="auth-screen">
    <form class="auth-card" on:submit|preventDefault={configurerMotsDePasse}>
      <img src={logo} alt="Gestion de Stock" />
      <h1>Créer les accès</h1>
      <p>Définissez deux mots de passe différents. Ils seront enregistrés de manière sécurisée.</p>
      {#if authError}<div class="auth-error">{authError}</div>{/if}
      <label for="new-admin-password">Mot de passe administrateur</label>
      <input id="new-admin-password" type="password" minlength="8" bind:value={adminPassword} autocomplete="new-password" required />
      <label for="confirm-admin-password">Confirmer le mot de passe administrateur</label>
      <input id="confirm-admin-password" type="password" minlength="8" bind:value={adminPasswordConfirm} autocomplete="new-password" required />
      <label for="new-user-password">Mot de passe utilisateur</label>
      <input id="new-user-password" type="password" minlength="8" bind:value={userPassword} autocomplete="new-password" required />
      <label for="confirm-user-password">Confirmer le mot de passe utilisateur</label>
      <input id="confirm-user-password" type="password" minlength="8" bind:value={userPasswordConfirm} autocomplete="new-password" required />
      <button type="submit" disabled={authBusy}>{authBusy ? 'Enregistrement...' : 'Enregistrer les mots de passe'}</button>
    </form>
  </div>
{:else if !currentRole}
  <div class="auth-screen">
    <form class="auth-card" on:submit|preventDefault={ouvrirSession}>
      <img src={logo} alt="Gestion de Stock" />
      <h1>Connexion</h1>
      <p>Choisissez votre profil puis saisissez votre mot de passe.</p>
      {#if authError}<div class="auth-error">{authError}</div>{/if}
      <label for="login-role">Profil</label>
      <select id="login-role" bind:value={selectedRole}>
        <option value="admin">Administrateur</option>
        <option value="user">Utilisateur</option>
      </select>
      <label for="login-password">Mot de passe</label>
      <input id="login-password" type="password" bind:value={password} autocomplete="current-password" required />
      <button type="submit" disabled={authBusy}>{authBusy ? 'Connexion...' : 'Se connecter'}</button>

      <div class="network-link-box">
        {#if serverInfo.mode === 'client'}
          <span class="net-client-badge">📡 Lié au PC 1</span>
          <button type="button" class="btn-net-link" on:click={() => modalReseau = true}>Changer d'IP</button>
        {:else}
          <button type="button" class="btn-net-link" on:click={() => modalReseau = true}>🔗 Connecter ce poste au PC 1 (Réseau)</button>
        {/if}
      </div>
    </form>
  </div>
{:else}
<div class="app-layout">
  <aside class="sidebar">
    <div class="brand">
      <img class="brand-logo" src={logo} alt="Logo Gestion de Stock" />
      <div class="brand-text">
        <span class="brand-title">Gestion</span>
        <span class="brand-sub">de Stock</span>
      </div>
    </div>

    <nav class="nav">
      <button
        class="nav-item {pageActive === 'caisse' ? 'active' : ''}"
        on:click={() => changerPage('caisse')}
      >
        <span class="icon">🛒</span>
        <span>Caisse</span>
      </button>

      {#if currentRole === 'admin'}
      <button
        class="nav-item {pageActive === 'produits' ? 'active' : ''}"
        on:click={() => changerPage('produits')}
      >
        <span class="icon">📦</span>
        <span>Produits</span>
      </button>

      <button
        class="nav-item {pageActive === 'stock' ? 'active' : ''}"
        on:click={() => changerPage('stock')}
      >
        <span class="icon">📊</span>
        <span>Stock</span>
      </button>
      {/if}

      <button
        class="nav-item {pageActive === 'factures' ? 'active' : ''}"
        on:click={() => changerPage('factures')}
      >
        <span class="icon">🧾</span>
        <span>Factures & Bons</span>
      </button>

      {#if currentRole === 'admin'}
      <button
        class="nav-item {pageActive === 'parametres' ? 'active' : ''}"
        on:click={() => changerPage('parametres')}
      >
        <span class="icon">⚙️</span>
        <span>Paramètres</span>
      </button>
      {/if}
    </nav>

    <div class="sidebar-footer">
      {#if serverInfo.running && serverInfo.ips.length > 0}
        <div class="serveur-badge" title="Le PC 2 peut se connecter à cette adresse IP">
          <span class="dot-green"></span>
          <span>PC 1 : {serverInfo.ips[0]}:8765</span>
        </div>
      {:else if serverInfo.mode === 'client'}
        <div class="serveur-badge client" title="Connecté au PC 1">
          <span class="dot-blue"></span>
          <span>Lié au PC 1</span>
        </div>
      {/if}
      <span class="role-label">{currentRole === 'admin' ? 'Administrateur' : 'Utilisateur'}</span>
      <button class="btn-logout" on:click={fermerSession}>Déconnexion</button>
      {#if currentRole === 'admin'}
        <button class="btn-apropos" on:click={() => modalAPropos = true}>
          ℹ️ À propos
        </button>
      {/if}
    </div>
  </aside>

  <main class="main-content">
    {#if pageActive === 'caisse'}
      <Caisse />
    {:else if pageActive === 'produits'}
      <Produits />
    {:else if pageActive === 'stock'}
      <Stock />
    {:else if pageActive === 'factures'}
      <Factures />
    {:else if pageActive === 'parametres'}
      <Parametres />
    {/if}
  </main>
</div>
{/if}

{#if modalReseau}
  <div class="overlay" on:click={() => modalReseau = false}>
    <div class="modal modal-reseau" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Liaison Réseau entre 2 Postes</h2>
        <button class="close-btn" on:click={() => modalReseau = false}>&times;</button>
      </div>
      <div class="modal-body">
        <p>Pour que le <strong>PC 2</strong> accède à la même base de données que le <strong>PC 1</strong>, indiquez l'adresse IP affichée sur le PC 1.</p>
        {#if erreurReseau}<div class="auth-error">{erreurReseau}</div>{/if}
        <div class="form-group-modal">
          <label for="ip-pc1">Adresse IP du PC 1 :</label>
          <input id="ip-pc1" type="text" placeholder="Ex: 192.168.1.50" bind:value={ipSaisie} />
          <small class="hint">Le PC 1 doit être allumé et l'administrateur doit y être connecté.</small>
        </div>
        <div class="btn-actions-modal">
          <button type="button" class="btn-primary-modal" disabled={reseauBusy} on:click={enregistrerLiaisonClient}>
            {reseauBusy ? 'Test de connexion...' : 'Valider la connexion au PC 1'}
          </button>
          {#if serverInfo.mode === 'client'}
            <button type="button" class="btn-secondary-modal" disabled={reseauBusy} on:click={repasserEnModeLocal}>
              Repasser en mode local (PC 1)
            </button>
          {/if}
        </div>
      </div>
    </div>
  </div>
{/if}

<APropos bind:visible={modalAPropos} onClose={() => modalAPropos = false} />

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
    color: #1e293b;
    background-color: #f8fafc;
    user-select: none;
    -webkit-user-select: none;
  }
  .auth-screen {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: #0f172a;
    padding: 24px;
    box-sizing: border-box;
  }
  .auth-card {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
    max-width: 440px;
    padding: 32px;
    border-radius: 14px;
    background: white;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.3);
    box-sizing: border-box;
  }
  .auth-card img {
    width: 64px;
    height: 64px;
    margin: 0 auto 4px auto;
    object-fit: contain;
  }
  .auth-card h1 {
    margin: 0;
    font-size: 1.45rem;
    text-align: center;
    color: #0f172a;
  }
  .auth-card p {
    margin: 0 0 8px 0;
    font-size: 0.9rem;
    color: #64748b;
    text-align: center;
  }
  .auth-card label {
    font-size: 0.85rem;
    color: #334155;
    font-weight: 600;
  }
  .auth-card input, .auth-card select {
    box-sizing: border-box;
    width: 100%;
    padding: 11px 12px;
    border: 1px solid #cbd5e1;
    border-radius: 7px;
    font: inherit;
  }
  .auth-card button {
    margin-top: 8px;
    padding: 12px;
    border: 0;
    border-radius: 7px;
    color: white;
    background: #0284c7;
    font: inherit;
    font-weight: 700;
    cursor: pointer;
  }
  .auth-card button:disabled {
    opacity: 0.6;
    cursor: wait;
  }
  .auth-error {
    padding: 10px;
    border-radius: 6px;
    color: #991b1b;
    background: #fee2e2;
    font-size: 0.88rem;
  }
  .network-link-box {
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid #e2e8f0;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    font-size: 0.85rem;
  }
  .net-client-badge {
    color: #0369a1;
    font-weight: 600;
  }
  .btn-net-link {
    background: none !important;
    color: #0284c7 !important;
    border: none !important;
    padding: 0 !important;
    margin: 0 !important;
    font-size: 0.85rem !important;
    font-weight: 500 !important;
    text-decoration: underline;
    cursor: pointer;
  }
  .app-layout {
    display: flex;
    height: 100vh;
    overflow: hidden;
  }
  .sidebar {
    width: 240px;
    background: #0f172a;
    color: white;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
  }
  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 24px 20px;
    border-bottom: 1px solid #1e293b;
  }
  .brand-logo {
    width: 52px;
    height: 52px;
    flex-shrink: 0;
    object-fit: contain;
  }
  .brand-text {
    display: flex;
    flex-direction: column;
  }
  .brand-title {
    font-size: 1.1rem;
    font-weight: 700;
    color: #ffffff;
    line-height: 1.2;
  }
  .brand-sub {
    font-size: 0.8rem;
    color: #94a3b8;
  }
  .nav {
    flex: 1;
    padding: 16px 12px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    border: none;
    background: transparent;
    color: #cbd5e1;
    font-size: 0.95rem;
    font-weight: 500;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    transition: background 0.15s, color 0.15s;
  }
  .nav-item:hover {
    background: #1e293b;
    color: #ffffff;
  }
  .nav-item.active {
    background: #0284c7;
    color: #ffffff;
    font-weight: 600;
  }
  .nav-item .icon {
    font-size: 1.2rem;
  }
  .sidebar-footer {
    padding: 16px;
    border-top: 1px solid #1e293b;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .serveur-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    background: #1e293b;
    border: 1px solid #334155;
    padding: 6px 10px;
    border-radius: 6px;
    font-size: 0.78rem;
    color: #cbd5e1;
  }
  .serveur-badge.client {
    border-color: #0284c7;
  }
  .dot-green {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #22c55e;
    flex-shrink: 0;
  }
  .dot-blue {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #38bdf8;
    flex-shrink: 0;
  }
  .role-label {
    color: #cbd5e1;
    text-align: center;
    font-size: 0.85rem;
  }
  .btn-logout {
    padding: 10px;
    color: white;
    border: 1px solid #475569;
    border-radius: 6px;
    background: transparent;
    cursor: pointer;
  }
  .btn-logout:hover {
    background: #334155;
  }
  .btn-apropos {
    width: 100%;
    padding: 10px 14px;
    background: #1e293b;
    border: 1px solid #334155;
    color: #cbd5e1;
    font-size: 0.85rem;
    border-radius: 6px;
    cursor: pointer;
    text-align: left;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: background 0.15s, color 0.15s;
  }
  .btn-apropos:hover {
    background: #334155;
    color: #ffffff;
  }
  .main-content {
    flex: 1;
    overflow-y: auto;
    background-color: #f8fafc;
  }
  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(15, 23, 42, 0.65);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .modal-reseau {
    background: white;
    border-radius: 12px;
    width: 100%;
    max-width: 480px;
    padding: 24px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.3);
  }
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }
  .modal-header h2 {
    margin: 0;
    font-size: 1.25rem;
    color: #0f172a;
  }
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: #64748b;
    cursor: pointer;
  }
  .form-group-modal {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin: 16px 0;
  }
  .form-group-modal label {
    font-weight: 600;
    font-size: 0.9rem;
    color: #334155;
  }
  .form-group-modal input {
    padding: 10px 12px;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 1rem;
  }
  .hint {
    color: #64748b;
    font-size: 0.8rem;
  }
  .btn-actions-modal {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-top: 20px;
  }
  .btn-primary-modal {
    padding: 12px;
    border: 0;
    border-radius: 7px;
    color: white;
    background: #0284c7;
    font-weight: 700;
    cursor: pointer;
  }
  .btn-secondary-modal {
    padding: 10px;
    border: 1px solid #cbd5e1;
    border-radius: 7px;
    color: #475569;
    background: white;
    font-weight: 600;
    cursor: pointer;
  }
</style>
