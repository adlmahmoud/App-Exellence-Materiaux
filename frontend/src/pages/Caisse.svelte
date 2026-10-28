<script>
  import { onMount } from 'svelte';
  import {
    GetProduitParCodeBarres,
    ListerProduits,
    CreerFacture,
    ImprimerFacturePDF,
    CreerBon,
    ImprimerBonPDF,
    LancerPaiementTPE,
    AnnulerPaiementTPE
  } from '../lib/api.js';

  let codeBarresSaisi = '';
  let panier = [];
  let messageErreur = '';
  let messageSucces = '';
  let inputScanElement;

  let modalTPE = false;
  let statutTPE = 'EN_COURS';
  let messageTPE = '';
  let detailTPE = null;

  let suggestions = [];
  let timerSuggestions;

  let cacheProduits = new Map();

  let azertyMap = {
    '&': '1', 'é': '2', '"': '3', "'": '4', '(': '5',
    '-': '6', 'è': '7', '_': '8', 'ç': '9', 'à': '0'
  };

  function normaliserCodeBarres(str) {
    if (!str) return '';
    let res = '';
    for (let i = 0; i < str.length; i++) {
      const ch = str[i];
      res += azertyMap[ch] !== undefined ? azertyMap[ch] : ch;
    }
    return res.trim();
  }

  function bipScan() {
    try {
      const AudioCtx = window.AudioContext || window.webkitAudioContext;
      if (!AudioCtx) return;
      const ctx = new AudioCtx();
      const osc = ctx.createOscillator();
      const gain = ctx.createGain();
      osc.type = 'sine';
      osc.frequency.setValueAtTime(880, ctx.currentTime);
      gain.gain.setValueAtTime(0.08, ctx.currentTime);
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.08);
      osc.connect(gain);
      gain.connect(ctx.destination);
      osc.start();
      osc.stop(ctx.currentTime + 0.08);
    } catch (e) {}
  }

  async function chargerCacheProduits() {
    try {
      const res = await ListerProduits({ limit: 500, actif: true });
      if (res && res.produits) {
        cacheProduits.clear();
        for (const p of res.produits) {
          if (p.codeBarres) {
            const cb = p.codeBarres.trim();
            cacheProduits.set(cb.toLowerCase(), p);
            cacheProduits.set(normaliserCodeBarres(cb).toLowerCase(), p);
          }
          if (p.reference) {
            cacheProduits.set(p.reference.trim().toLowerCase(), p);
          }
        }
      }
    } catch (e) {}
  }

  let dernierInputTime = 0;
  let inputBurstCount = 0;
  let timerBurstInput;

  let globalScanBuffer = '';
  let globalLastKeyTime = 0;
  let timerBurstGlobal;

  let clientNom = 'Client Comptoir';
  let modeReglement = 'especes';
  let chargementPaiement = false;

  $: totalHT = panier.reduce((acc, item) => acc + (item.prixVenteHt * item.quantite), 0);
  $: totalTTC = panier.reduce((acc, item) => acc + (item.prixVenteTtc * item.quantite), 0);
  $: totalTVA = totalTTC - totalHT;

  onMount(() => {
    chargerCacheProduits();
    focusScanInput();
    const handleWindowFocus = () => focusScanInput();
    document.addEventListener('keydown', capturerScanGlobal, true);
    window.addEventListener('focus', handleWindowFocus);
    return () => {
      document.removeEventListener('keydown', capturerScanGlobal, true);
      window.removeEventListener('focus', handleWindowFocus);
      clearTimeout(timerSuggestions);
      clearTimeout(timerBurstInput);
      clearTimeout(timerBurstGlobal);
    };
  });

  function focusScanInput() {
    inputScanElement?.focus();
  }

  function onInputScan() {
    const maintenant = performance.now();
    const intervalle = maintenant - dernierInputTime;
    dernierInputTime = maintenant;

    if (intervalle < 65) {
      inputBurstCount++;
    } else {
      inputBurstCount = 1;
    }

    clearTimeout(timerBurstInput);
    clearTimeout(timerSuggestions);

    if (inputBurstCount >= 3) {
      timerBurstInput = setTimeout(() => {
        const code = codeBarresSaisi.trim();
        if (code.length >= 2) {
          void traiterCodeBarres(code);
        }
        inputBurstCount = 0;
      }, 70);
      return;
    }

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
    if (e.key === 'Escape') {
      suggestions = [];
      return;
    }
    if (e.key !== 'Enter' && e.key !== 'Tab') return;

    e.preventDefault();
    if (e.repeat) return;
    clearTimeout(timerBurstInput);
    const code = codeBarresSaisi.trim();
    if (!code) return;
    clearTimeout(timerSuggestions);
    await traiterCodeBarres(code);
  }

  function capturerScanGlobal(e) {
    if (e.ctrlKey || e.altKey || e.metaKey) return;
    if (['Shift', 'CapsLock', 'Control', 'Alt', 'Meta', 'AltGraph'].includes(e.key)) return;
    if (document.activeElement === inputScanElement) return;

    const maintenant = performance.now();
    const intervalle = maintenant - globalLastKeyTime;
    globalLastKeyTime = maintenant;

    if (e.key === 'Enter' || e.key === 'Tab') {
      clearTimeout(timerBurstGlobal);
      if (globalScanBuffer.length >= 2) {
        e.preventDefault();
        e.stopPropagation();
        const code = globalScanBuffer;
        globalScanBuffer = '';
        void traiterCodeBarres(code);
      }
      return;
    }

    if (e.key.length !== 1) return;

    if (intervalle > 200) {
      globalScanBuffer = '';
    }

    globalScanBuffer += e.key;

    clearTimeout(timerBurstGlobal);
    if (globalScanBuffer.length >= 3) {
      timerBurstGlobal = setTimeout(() => {
        const code = globalScanBuffer;
        globalScanBuffer = '';
        void traiterCodeBarres(code);
      }, 70);
    }
  }

  async function trouverProduitParCode(codeBrut) {
    if (!codeBrut) return null;
    const codeNet = codeBrut.trim();
    const codeNorm = normaliserCodeBarres(codeNet);

    const cleNet = codeNet.toLowerCase();
    const cleNorm = codeNorm.toLowerCase();

    if (cacheProduits.has(cleNet)) return cacheProduits.get(cleNet);
    if (cacheProduits.has(cleNorm)) return cacheProduits.get(cleNorm);

    try {
      const prod = await GetProduitParCodeBarres(codeNorm);
      if (prod && prod.id) {
        cacheProduits.set(cleNorm, prod);
        cacheProduits.set(cleNet, prod);
        return prod;
      }
    } catch (e) {}

    if (codeNorm !== codeNet) {
      try {
        const prod = await GetProduitParCodeBarres(codeNet);
        if (prod && prod.id) {
          cacheProduits.set(cleNorm, prod);
          cacheProduits.set(cleNet, prod);
          return prod;
        }
      } catch (e) {}
    }

    try {
      const res = await ListerProduits({ recherche: codeNorm, limit: 1 });
      if (res && res.produits && res.produits.length > 0) {
        const prod = res.produits[0];
        cacheProduits.set(cleNorm, prod);
        cacheProduits.set(cleNet, prod);
        return prod;
      }
    } catch (e) {}

    return null;
  }

  async function traiterCodeBarres(code) {
    if (!code) return;
    clearTimeout(timerSuggestions);
    clearTimeout(timerBurstInput);
    clearTimeout(timerBurstGlobal);
    messageErreur = '';
    messageSucces = '';
    suggestions = [];
    codeBarresSaisi = '';

    const prod = await trouverProduitParCode(code);
    if (prod) {
      ajouterAuPanier(prod);
      bipScan();
    } else {
      messageErreur = 'Aucun produit trouvé pour : ' + code;
    }

    focusScanInput();
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
    if (modeReglement === 'cb') {
      await demarrerPaiementTPE();
      return;
    }
    await finaliserFacturation();
  }

  async function demarrerPaiementTPE() {
    statutTPE = 'EN_COURS';
    messageTPE = 'Envoi du montant au terminal de paiement...';
    detailTPE = null;
    modalTPE = true;

    try {
      const montantEuros = totalTTC / 100;
      const res = await LancerPaiementTPE(montantEuros, 'VENTE');
      detailTPE = res;

      if (res.statut === 'VALIDE') {
        statutTPE = 'VALIDE';
        messageTPE = res.message || 'Paiement accepté !';
        setTimeout(async () => {
          modalTPE = false;
          await finaliserFacturation();
        }, 1200);
      } else if (res.statut === 'REFUSE') {
        statutTPE = 'REFUSE';
        messageTPE = res.message || 'Paiement refusé par la banque';
      } else {
        statutTPE = 'ANNULE';
        messageTPE = res.message || 'Transaction annulée sur le terminal';
      }
    } catch (err) {
      statutTPE = 'ANNULE';
      messageTPE = 'Erreur terminal : ' + err;
    }
  }

  async function annulerTPEEtFermer() {
    try {
      await AnnulerPaiementTPE();
    } catch (e) {}
    modalTPE = false;
  }

  function changerModeEtFermer(nouveauMode) {
    modalTPE = false;
    modeReglement = nouveauMode;
  }

  async function finaliserFacturation() {
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
        on:click={focusScanInput}
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

  {#if modalTPE}
    <div class="modal-tpe-overlay">
      <div class="modal-tpe-card">
        <div class="modal-tpe-header">
          <h3>💳 Paiement par Carte Bancaire (TPE)</h3>
        </div>
        <div class="modal-tpe-body">
          <div class="tpe-montant-box">
            <span class="tpe-montant-label">Montant envoyé au terminal :</span>
            <span class="tpe-montant-val">{(totalTTC / 100).toFixed(2)} €</span>
          </div>

          {#if statutTPE === 'EN_COURS'}
            <div class="tpe-box en-cours">
              <div class="tpe-spinner"></div>
              <h4>Transaction en cours sur le TPE...</h4>
              <p>Le client doit insérer sa carte ou utiliser le sans contact sur le terminal.</p>
              <span class="tpe-msg">{messageTPE}</span>
            </div>
            <button class="btn-tpe-cancel" on:click={annulerTPEEtFermer}>
              Annuler depuis la caisse
            </button>
          {:else if statutTPE === 'VALIDE'}
            <div class="tpe-box valide">
              <div class="tpe-icon">✅</div>
              <h4>Paiement Validé !</h4>
              <p>{messageTPE}</p>
              {#if detailTPE && detailTPE.numeroAutorisation}
                <div class="tpe-badge-auth">Autorisation : {detailTPE.numeroAutorisation}</div>
              {/if}
              <span class="tpe-msg">Validation et création de la facture en cours...</span>
            </div>
          {:else if statutTPE === 'REFUSE'}
            <div class="tpe-box refuse">
              <div class="tpe-icon">❌</div>
              <h4>Transaction Refusée par la banque</h4>
              <p>{messageTPE}</p>
              <small>Motif : Solde insuffisant, code PIN erroné ou plafond dépassé</small>
            </div>
            <div class="tpe-actions">
              <button class="btn-tpe-retry" on:click={demarrerPaiementTPE}>
                🔄 Recommencer sur le TPE
              </button>
              <button class="btn-tpe-cash" on:click={() => changerModeEtFermer('especes')}>
                💵 Payer en Espèces
              </button>
              <button class="btn-tpe-close" on:click={() => modalTPE = false}>
                Fermer
              </button>
            </div>
          {:else}
            <div class="tpe-box annule">
              <div class="tpe-icon">⚠️</div>
              <h4>Transaction Annulée</h4>
              <p>{messageTPE}</p>
              <small>Annulation par le client sur le TPE ou abandon</small>
            </div>
            <div class="tpe-actions">
              <button class="btn-tpe-retry" on:click={demarrerPaiementTPE}>
                🔄 Réessayer le paiement CB
              </button>
              <button class="btn-tpe-cash" on:click={() => changerModeEtFermer('especes')}>
                💵 Payer en Espèces
              </button>
              <button class="btn-tpe-close" on:click={() => modalTPE = false}>
                Fermer
              </button>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
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
  .modal-tpe-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(15, 23, 42, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    backdrop-filter: blur(3px);
  }
  .modal-tpe-card {
    background: #ffffff;
    width: 90%;
    max-width: 480px;
    border-radius: 12px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2);
    overflow: hidden;
  }
  .modal-tpe-header {
    background: #0f172a;
    color: white;
    padding: 16px 20px;
  }
  .modal-tpe-header h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 700;
  }
  .modal-tpe-body {
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .tpe-montant-box {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 14px 18px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .tpe-montant-label {
    font-size: 0.95rem;
    color: #64748b;
    font-weight: 600;
  }
  .tpe-montant-val {
    font-size: 1.5rem;
    font-weight: 800;
    color: #0284c7;
  }
  .tpe-box {
    border-radius: 10px;
    padding: 20px;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }
  .tpe-box.en-cours {
    background: #f0f9ff;
    border: 2px solid #bae6fd;
    color: #0369a1;
  }
  .tpe-box.valide {
    background: #ecfdf5;
    border: 2px solid #a7f3d0;
    color: #065f46;
  }
  .tpe-box.refuse {
    background: #fef2f2;
    border: 2px solid #fecaca;
    color: #991b1b;
  }
  .tpe-box.annule {
    background: #fffbeb;
    border: 2px solid #fde68a;
    color: #92400e;
  }
  .tpe-spinner {
    width: 44px;
    height: 44px;
    border: 4px solid #bae6fd;
    border-top-color: #0284c7;
    border-radius: 50%;
    animation: spinTPE 0.8s linear infinite;
    margin-bottom: 6px;
  }
  @keyframes spinTPE {
    to { transform: rotate(360deg); }
  }
  .tpe-icon {
    font-size: 2.5rem;
    line-height: 1;
  }
  .tpe-box h4 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 700;
  }
  .tpe-box p {
    margin: 0;
    font-size: 0.9rem;
  }
  .tpe-msg {
    font-size: 0.8rem;
    opacity: 0.85;
  }
  .tpe-badge-auth {
    background: #059669;
    color: white;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 700;
  }
  .btn-tpe-cancel {
    background: #fee2e2;
    color: #dc2626;
    border: 1px solid #fca5a5;
    padding: 10px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    width: 100%;
  }
  .btn-tpe-cancel:hover {
    background: #fecaca;
  }
  .tpe-actions {
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 100%;
  }
  .btn-tpe-retry {
    background: #16a34a;
    color: white;
    border: none;
    padding: 12px;
    border-radius: 6px;
    font-weight: 700;
    cursor: pointer;
    font-size: 0.95rem;
  }
  .btn-tpe-retry:hover {
    background: #15803d;
  }
  .btn-tpe-cash {
    background: #0284c7;
    color: white;
    border: none;
    padding: 10px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-tpe-cash:hover {
    background: #0369a1;
  }
  .btn-tpe-close {
    background: #e2e8f0;
    color: #475569;
    border: none;
    padding: 10px;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
  }
  .btn-tpe-close:hover {
    background: #cbd5e1;
  }
</style>
