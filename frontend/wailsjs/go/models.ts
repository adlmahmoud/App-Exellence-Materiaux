export namespace parametres {
	
	export class ParametresEntreprise {
	    id: string;
	    denomination: string;
	    formeJuridique: string;
	    capital: number;
	    adresse: string;
	    codePostal: string;
	    ville: string;
	    siren: string;
	    siret: string;
	    rcs: string;
	    numeroTva: string;
	    telephone: string;
	    email: string;
	    siteWeb: string;
	    logoPath: string;
	    devise: string;
	    tauxTvaDefaut: number;
	    tauxTvaDisponibles: string;
	    mentionsLegales: string;
	    conditionsReglement: string;
	    formatNumerotationFacture: string;
	    formatNumerotationBon: string;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ParametresEntreprise(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.denomination = source["denomination"];
	        this.formeJuridique = source["formeJuridique"];
	        this.capital = source["capital"];
	        this.adresse = source["adresse"];
	        this.codePostal = source["codePostal"];
	        this.ville = source["ville"];
	        this.siren = source["siren"];
	        this.siret = source["siret"];
	        this.rcs = source["rcs"];
	        this.numeroTva = source["numeroTva"];
	        this.telephone = source["telephone"];
	        this.email = source["email"];
	        this.siteWeb = source["siteWeb"];
	        this.logoPath = source["logoPath"];
	        this.devise = source["devise"];
	        this.tauxTvaDefaut = source["tauxTvaDefaut"];
	        this.tauxTvaDisponibles = source["tauxTvaDisponibles"];
	        this.mentionsLegales = source["mentionsLegales"];
	        this.conditionsReglement = source["conditionsReglement"];
	        this.formatNumerotationFacture = source["formatNumerotationFacture"];
	        this.formatNumerotationBon = source["formatNumerotationBon"];
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace produits {
	
	export class FiltreProduits {
	    recherche: string;
	    alerteStock: boolean;
	    actif?: boolean;
	    page: number;
	    limit: number;
	    triChamp: string;
	    triAscendant: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FiltreProduits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recherche = source["recherche"];
	        this.alerteStock = source["alerteStock"];
	        this.actif = source["actif"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.triChamp = source["triChamp"];
	        this.triAscendant = source["triAscendant"];
	    }
	}
	export class ModificationProduit {
	    id: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    description: string;
	    prixAchatHt: number;
	    prixVenteHt: number;
	    prixVenteTtc: number;
	    tauxTva: number;
	    stockMinimum: number;
	    unite: string;
	    actif: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModificationProduit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.description = source["description"];
	        this.prixAchatHt = source["prixAchatHt"];
	        this.prixVenteHt = source["prixVenteHt"];
	        this.prixVenteTtc = source["prixVenteTtc"];
	        this.tauxTva = source["tauxTva"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	        this.actif = source["actif"];
	    }
	}
	export class NouveauProduit {
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    description: string;
	    prixAchatHt: number;
	    prixVenteHt: number;
	    prixVenteTtc: number;
	    tauxTva: number;
	    stockInitial: number;
	    stockMinimum: number;
	    unite: string;
	
	    static createFrom(source: any = {}) {
	        return new NouveauProduit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.description = source["description"];
	        this.prixAchatHt = source["prixAchatHt"];
	        this.prixVenteHt = source["prixVenteHt"];
	        this.prixVenteTtc = source["prixVenteTtc"];
	        this.tauxTva = source["tauxTva"];
	        this.stockInitial = source["stockInitial"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	    }
	}
	export class Produit {
	    id: string;
	    codeBarres: string;
	    reference: string;
	    designation: string;
	    description: string;
	    prixAchatHt: number;
	    prixVenteHt: number;
	    prixVenteTtc: number;
	    tauxTva: number;
	    stockActuel: number;
	    stockMinimum: number;
	    unite: string;
	    actif: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Produit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.codeBarres = source["codeBarres"];
	        this.reference = source["reference"];
	        this.designation = source["designation"];
	        this.description = source["description"];
	        this.prixAchatHt = source["prixAchatHt"];
	        this.prixVenteHt = source["prixVenteHt"];
	        this.prixVenteTtc = source["prixVenteTtc"];
	        this.tauxTva = source["tauxTva"];
	        this.stockActuel = source["stockActuel"];
	        this.stockMinimum = source["stockMinimum"];
	        this.unite = source["unite"];
	        this.actif = source["actif"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PageProduits {
	    produits: Produit[];
	    total: number;
	    page: number;
	    limit: number;
	    totalPages: number;
	
	    static createFrom(source: any = {}) {
	        return new PageProduits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.produits = this.convertValues(source["produits"], Produit);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.limit = source["limit"];
	        this.totalPages = source["totalPages"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

