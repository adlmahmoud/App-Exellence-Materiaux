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

