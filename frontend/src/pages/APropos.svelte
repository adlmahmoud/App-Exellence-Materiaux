<script>
  import logo from '../assets/images/gestion-stock.png';
  import { GetAPropos } from '../../wailsjs/go/main/App.js';
  import { onMount } from 'svelte';

  export let visible = false;
  export let onClose = () => {};

  let info = {
    nom: 'Gestion de Stock',
    version: '1.0.0',
    annee: '2026',
    auteurs: 'Adel Mahmoud et Sidou'
  };

  onMount(async () => {
    try {
      const res = await GetAPropos();
      if (res) info = res;
    } catch (e) {}
  });
</script>

{#if visible}
  <div class="overlay" on:click={onClose}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>À propos</h2>
        <button class="close-btn" on:click={onClose}>&times;</button>
      </div>
      <div class="modal-body">
        <img class="app-badge" src={logo} alt="Logo Gestion de Stock" />
        <h3>{info.nom}</h3>
        <p class="subtitle">Logiciel de gestion de stock et facturation</p>
        <div class="details">
          <div class="detail-row">
            <span class="label">Version :</span>
            <span class="val">{info.version}</span>
          </div>
          <div class="detail-row">
            <span class="label">Année :</span>
            <span class="val">{info.annee}</span>
          </div>
          <div class="detail-row">
            <span class="label">Auteurs :</span>
            <span class="val">{info.auteurs}</span>
          </div>
          <div class="detail-row">
            <span class="label">Licence :</span>
            <span class="val">Tous droits réservés</span>
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-primary" on:click={onClose}>Fermer</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(15, 23, 42, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .modal {
    background: #ffffff;
    border-radius: 12px;
    width: 420px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2);
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
    font-size: 1.1rem;
    color: #1e293b;
  }
  .close-btn {
    background: transparent;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #64748b;
  }
  .modal-body {
    padding: 24px 20px;
    text-align: center;
  }
  .app-badge {
    width: 112px;
    height: 112px;
    object-fit: contain;
    margin: 0 auto 16px;
  }
  .modal-body h3 {
    margin: 0 0 4px;
    color: #0f172a;
    font-size: 1.3rem;
  }
  .subtitle {
    margin: 0 0 20px;
    color: #64748b;
    font-size: 0.9rem;
  }
  .details {
    background: #f1f5f9;
    border-radius: 8px;
    padding: 12px 16px;
    text-align: left;
  }
  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 6px 0;
    font-size: 0.9rem;
    border-bottom: 1px solid #e2e8f0;
  }
  .detail-row:last-child {
    border-bottom: none;
  }
  .label {
    color: #64748b;
  }
  .val {
    font-weight: 600;
    color: #1e293b;
  }
  .modal-footer {
    padding: 12px 20px;
    background: #f8fafc;
    border-top: 1px solid #e2e8f0;
    text-align: right;
  }
  .btn-primary {
    background: #0284c7;
    color: white;
    border: none;
    padding: 8px 18px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 500;
  }
</style>
