# Exellence Matériaux

Logiciel de bureau tout-en-un : Caisse & scan rapide, gestion des stocks, facturation et bons de livraison pour **Exellence Matériaux**.  
Conçu et développé par **Adel Mahmoud et Sidou**.

---

## ⚡ Prérequis (Pour l'utilisateur)

**AUCUNE CONNAISSANCE EN CODE N'EST NÉCESSAIRE.**  
Le logiciel est **déjà compilé et 100% prêt à l'emploi**.

- ❌ **Pas besoin** d'installer Go.
- ❌ **Pas besoin** d'installer Node.js ni de terminal.
- ❌ **Pas de lignes de commande** à taper.
- ✔️ **Le seul prérequis** : Un ordinateur sous **Windows 10 ou Windows 11**.

*(Note : Dans 99% des cas, Windows possède déjà le composant d'affichage WebView2. S'il venait à manquer sur un vieux PC, il suffit de [cliquer ici pour l'installer en 30 secondes](https://go.microsoft.com/fwlink/p/?LinkId=2124703)).*

---

## 🚀 Installation & Lancement en 2 étapes

### 1. Extraire le dossier
Si vous avez reçu l'archive **`Exellence-Materiaux-Complet.zip`** :
1. Faites un **Clic droit** sur le fichier ZIP.
2. Cliquez sur **« Extraire tout... »** puis validez (par exemple sur votre **Bureau**).
> ⚠️ **Important** : Ne lancez pas le logiciel directement depuis le fichier ZIP sans l'extraire !

### 2. Démarrer le logiciel
- Ouvrez le dossier extrait et double-cliquez directement sur **`Exellence-Materiaux.exe`** (ou sur **`lancer.bat`**).
- Le logiciel s'ouvre directement à l'écran !

---

## 🛠️ Que faire si Windows affiche un message au premier lancement ?

### 1. Écran bleu « Windows a protégé votre ordinateur » (SmartScreen)
C'est une alerte de sécurité standard de Windows lorsque vous recevez un nouveau fichier `.exe` par WhatsApp, Discord, WeTransfer ou clé USB.

Pour l'ouvrir en 2 clics :
1. Cliquez sur le lien **« Informations complémentaires »** sous le message.
2. Cliquez sur le bouton **« Exécuter quand même »** qui apparaît.  
*(L'application se lance et ne vous le redemandera plus).*

*Astuce alternative :* Clic droit sur `Exellence-Materiaux.exe` ➔ **Propriétés** ➔ tout en bas, cochez la case **« Débloquer »** ➔ **Appliquer** ➔ **OK**.

---

## 📦 Ce que fait le logiciel

- **🛒 Caisse & Douchette USB** : Branchez votre lecteur code-barres USB, scannez vos articles (ou tapez le code + Entrée), encaissez et imprimez la facture directement.
- **📦 Produits & Stocks** : Catalogue avec recherche instantanée, entrées de stock fournisseurs, sorties, casse, ajustements physiques d'inventaire et alertes automatiques de stock bas.
- **🧾 Factures & Bons de livraison** : Numérotation séquentielle officielle sans trou (`FAC-2026-XXXXX`), modèle conforme avec logo officiel haute définition, calcul automatique HT/TVA/TTC et génération de PDF instantanée.
- **⚙️ Paramètres & Sauvegarde** : Sauvegarde sécurisée de toutes vos données en un clic.

---

## 🔒 Confidentialité des données de l'entreprise (Kbis)

Ce projet est configuré pour être mis sur **GitHub sans aucun risque** de fuite de données :
- Vos informations confidentielles (SIREN, RCS, adresse exacte) se trouvent dans le fichier local **`kbis.json`**.
- Ce fichier est strictement **bloqué et ignoré par Git** (il n'ira jamais sur GitHub).
- Pour transmettre le logiciel à un ami en privé, envoyez-lui directement le pack ZIP complet contenant le fichier `kbis.json`.

---

## 💻 Annexe : Pour les développeurs uniquement (Optionnel)

<details>
<summary>Cliquez ici pour afficher les commandes de compilation depuis les sources</summary>

Si vous souhaitez recompiler l'application vous-même depuis le code source :

**Prérequis développeur :**
- Go >= 1.22
- Node.js >= 18 & npm
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

**Commandes :**
```bash
# Lancement en mode développement
wails dev

# Compilation de l'exécutable final Windows
wails build -tags native_webview2loader
```
</details>

---

## 📄 Licence

Copyright © 2026 **Adel Mahmoud et Sidou**. Tous droits réservés.
