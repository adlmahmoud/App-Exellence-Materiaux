import * as localAPI from '../../wailsjs/go/main/App.js';

let networkConfigPromise;
let networkSessionToken = '';

function getNetworkConfig() {
  if (!networkConfigPromise) {
    networkConfigPromise = window.go.main.App.GetNetworkConfig();
  }
  return networkConfigPromise;
}

function chooseLogo() {
  return new Promise((resolve, reject) => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/png,image/jpeg';
    input.oncancel = () => resolve('');
    input.onchange = () => {
      const file = input.files?.[0];
      if (!file) {
        resolve('');
        return;
      }
      if (file.size > 10 * 1024 * 1024) {
        reject(new Error('Le logo ne doit pas dépasser 10 Mo.'));
        return;
      }
      const reader = new FileReader();
      reader.onerror = () => reject(new Error('Impossible de lire le fichier logo.'));
      reader.onload = () => resolve(String(reader.result).split(',')[1] || '');
      reader.readAsDataURL(file);
    };
    input.click();
  });
}

function downloadFile(result) {
  const bytes = Uint8Array.from(atob(result.data), character => character.charCodeAt(0));
  const type = result.filename.toLowerCase().endsWith('.pdf')
    ? 'application/pdf'
    : 'application/octet-stream';
  const url = URL.createObjectURL(new Blob([bytes], { type }));
  const link = document.createElement('a');
  link.href = url;
  link.download = result.filename;
  link.click();
  setTimeout(() => URL.revokeObjectURL(url), 60_000);
  return result.filename;
}

async function invoke(method, ...args) {
  if (method === 'SetNetworkMode' || method === 'GetServerInfo' || method.includes('TPE')) {
    if (method === 'SetNetworkMode') {
      networkConfigPromise = null;
    }
    return localAPI[method](...args);
  }

  const config = await getNetworkConfig();
  if (config.mode !== 'client') {
    return localAPI[method](...args);
  }

  if (method === 'GetCurrentRole' && !networkSessionToken) return '';
  if (!networkSessionToken && !['GetAuthStatus', 'SetupAccessPasswords', 'Login'].includes(method)) {
    throw new Error('Veuillez vous connecter pour continuer.');
  }

  if (method === 'ImporterLogo') {
    const imageData = await chooseLogo();
    if (!imageData) return '';
    args = [imageData];
  }

  const response = await fetch(`${config.serverUrl.replace(/\/+$/, '')}/api/rpc`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${config.token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ method, args, sessionToken: networkSessionToken })
  });
  const payload = await response.json();
  if (!response.ok) {
    throw new Error(payload.error || `Erreur serveur : ${response.status}`);
  }

  const result = payload.result;
  if (method === 'Login') networkSessionToken = result;
  if (method === 'Logout') networkSessionToken = '';
  if (method === 'ImprimerFacturePDF' || method === 'ImprimerBonPDF' || method === 'SauvegarderBase') {
    return downloadFile(result);
  }
  return result;
}

const methods = [
  'AjusterStock',
  'ChangeAccessPasswords',
  'CreerBon',
  'CreerFacture',
  'CreerMouvementStock',
  'CreerProduit',
  'GetAPropos',
  'GetAuthStatus',
  'GetBon',
  'GetCurrentRole',
  'GetFacture',
  'GetNetworkConfig',
  'GetParametres',
  'GetProduitParCodeBarres',
  'GetServerInfo',
  'ImporterLogo',
  'ImprimerBonPDF',
  'ImprimerFacturePDF',
  'ListerBons',
  'ListerFactures',
  'ListerMouvementsStock',
  'ListerProduits',
  'Login',
  'Logout',
  'ModifierProduit',
  'ObtenirAlertesStock',
  'SauvegarderBase',
  'SetNetworkMode',
  'SetupAccessPasswords',
  'SupprimerMouvementStock',
  'SupprimerProduit',
  'UpdateParametres',
  'LancerPaiementTPE',
  'AnnulerPaiementTPE',
  'TesterConnexionTPE',
  'GetConfigTPE',
  'SauvegarderConfigTPE',
  'ListerPortsCOMTPE',
  'SetSimulateurScenarioTPE'
];

export const API = Object.fromEntries(methods.map(method => [method, (...args) => invoke(method, ...args)]));
export const {
  AjusterStock,
  ChangeAccessPasswords,
  CreerBon,
  CreerFacture,
  CreerMouvementStock,
  CreerProduit,
  GetAPropos,
  GetAuthStatus,
  GetBon,
  GetCurrentRole,
  GetFacture,
  GetNetworkConfig,
  GetParametres,
  GetProduitParCodeBarres,
  GetServerInfo,
  ImporterLogo,
  ImprimerBonPDF,
  ImprimerFacturePDF,
  ListerBons,
  ListerFactures,
  ListerMouvementsStock,
  ListerProduits,
  Login,
  Logout,
  ModifierProduit,
  ObtenirAlertesStock,
  SauvegarderBase,
  SetNetworkMode,
  SetupAccessPasswords,
  SupprimerMouvementStock,
  SupprimerProduit,
  UpdateParametres,
  LancerPaiementTPE,
  AnnulerPaiementTPE,
  TesterConnexionTPE,
  GetConfigTPE,
  SauvegarderConfigTPE,
  ListerPortsCOMTPE,
  SetSimulateurScenarioTPE
} = API;
