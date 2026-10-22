# Gestion de Stock

Logiciel de bureau tout-en-un : Caisse & scan rapide, gestion des stocks, facturation et bons de livraison.  
Conçu et développé par **Adel Mahmoud et Sidou**.

---

## ⚡ Prérequis (Pour l'utilisateur)

**AUCUNE CONNAISSANCE EN CODE N'EST NÉCESSAIRE.**  
Le logiciel est **déjà compilé et 100% prêt à l'emploi**.

- ❌ **Pas besoin** d'installer Go.
- ❌ **Pas besoin** d'installer Node.js ni de terminal.
- ❌ **Pas de lignes de commande** à taper.
- ✔️ **Le seul prérequis** : Un ordinateur sous **Windows 10 ou Windows 11**. L'installateur ajoute WebView2 si nécessaire.

---

## 🚀 Installation

1. [Téléchargez `Gestion-de-Stock-Setup.exe`](./Gestion-de-Stock-Setup.exe) puis lancez-le.
2. Suivez les étapes de l'installateur. L'installation est limitée à votre compte Windows et ne nécessite pas les droits administrateur.
3. Lancez **Gestion de Stock** depuis le raccourci créé sur le Bureau ou dans le menu Démarrer.

La base de données, les logos et les PDF de l'installation sont enregistrés dans votre dossier personnel Windows (`%APPDATA%\Gestion de Stock`) et ne sont pas inclus dans l'installateur. Ils restent ainsi séparés des fichiers du programme.

---

## 🛠️ Que faire si Windows affiche un message au premier lancement ?

### 1. Écran bleu « Windows a protégé votre ordinateur » (SmartScreen)
C'est une alerte de sécurité standard de Windows lorsque vous recevez un nouveau fichier `.exe`.

Pour l'ouvrir en 2 clics :
1. Cliquez sur le lien **« Informations complémentaires »** sous le message.
2. Cliquez sur le bouton **« Exécuter quand même »** qui apparaît.  
*(L'application se lance et ne vous le redemandera plus).*

*Astuce alternative :* Clic droit sur `Gestion-de-Stock-Setup.exe` ➔ **Propriétés** ➔ tout en bas, cochez la case **« Débloquer »** ➔ **Appliquer** ➔ **OK**.

---

## 📦 Ce que fait le logiciel

- **🛒 Caisse & Douchette USB** : Branchez votre lecteur code-barres USB, scannez vos articles (ou tapez le code + Entrée), encaissez et imprimez la facture directement.
- **📦 Produits & Stocks** : Catalogue avec recherche instantanée, entrées de stock fournisseurs, sorties, casse, ajustements physiques d'inventaire et alertes automatiques de stock bas.
- **🧾 Factures & Bons de livraison** : Numérotation séquentielle officielle sans trou (`FAC-2026-XXXXX`), logo personnalisable, calcul automatique HT/TVA/TTC et génération de PDF instantanée.
- **⚙️ Paramètres & Sauvegarde** : Sauvegarde sécurisée de toutes vos données en un clic.

---

## 🔒 Confidentialité des données de l'entreprise (Kbis)

L'installateur **`Gestion-de-Stock-Setup.exe`** contient le programme, pas la base de données locale ni les logos et PDF personnels.
- Les données de chaque utilisateur sont conservées séparément dans `%APPDATA%\Gestion de Stock`.
- Pour distribuer l'application, partagez uniquement le fichier d'installation. N'incluez pas votre base `exellence.db`, vos fichiers `kbis.json`, vos sauvegardes ni vos documents PDF.
- Avant de publier les sources ou d'autres fichiers sur Internet, vérifiez qu'ils ne contiennent aucune donnée confidentielle.

---

## 💻 Annexe : Pour les développeurs uniquement (Optionnel)

<details>
<summary>Cliquez ici pour afficher les commandes de compilation depuis les sources</summary>

Si vous souhaitez recompiler l'application vous-même depuis le code source :

**Prérequis développeur :**
- Go >= 1.22
- Node.js >= 18 & npm
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- NSIS (pour générer l'installateur)

**Commandes :**
```bash
# Lancement en mode développement
wails dev

# Génération de l'installateur Windows, sans droits administrateur à l'installation
wails build -nsis -installscope user
```
</details>

---

## 📄 Licence

Copyright © 2026 **Adel Mahmoud et Sidou**. Tous droits réservés.
