---
title: Recipes
weight: 3300
pre: "<i class=\"fa fa-shopping-cart\" aria-hidden=\"true\"></i> "
---

### Recipes

A recipe is a pre-configured Mashling JSON file which can be customized or used as is for a specific gateway use case. The precompiled recipe binaries along with the JSON files are available at https://www.mashling.io/moreRecipes.

#### Creating a recipe
A [recipe](https://github.com/TIBCOSoftware/mashling-recipes) can be created by the Mashling CLI tool <should redirect to Mashling CLI topic> or by customizing an existing recipe in https://www.mashling.io/moreRecipes.

#### Adding a recipe
A recipe should be contained in its own folder under 'recipes' folder. The recipe folder is composed of a gateway json file, README.md, an optional icon image file, optional Gopkg.lock and Gopkg.toml files. In the absence of the icon image file, the default Mashling icon image is used by mashling.io for the recipe. When the icon image file is present, the Mashling json file should have an icon image file field as follows:

```
{
	"mashling_schema": "0.2",
	"gateway": {
		"name": "allRecipe",
		"version": "1.0.0",
		"display_name":"KafkaTrigger to KafkaPublisher",
		"display_image":"displayImage.svg",
		"..."
  }
}
```

If "display_name" field is present in the json, its value is used as the recipe name in mashling.io. Otherwise, the value of "name" field is used.

Gopkg.lock and Gopkg.toml contains the specific dependent library versions to be used during the recipe binary compilation. In the abasence of those, the latest versions of the dependent libraries are to be used. For more information about the dependencies versioning, refer to the [Mashling CLI documentation](https://github.com/TIBCOSoftware/mashling/blob/master/cli/docs/gateway.md).

#### Publishing a Recipe

[recipe_registry.json](https://github.com/TIBCOSoftware/mashling-recipes/blob/master/recipe_registry.json) contains the list of recipe providers and the recipes to publish. The recipe folder name should be added to the "publish" field for the recipe to be made available in mashling.io. For example, "event-dispatcher-router-mashling" and "rest-conditional-gateway" recipes binaries are built and made downloadable from mashling.io given the following recipe_registry.json. Setting "featured" to "true" adds the recipe to the featured recipe list in mashling.io.

```
{  
  "recipe_repos":[  
    {  
      "provider":"TIBCOSoftware Engineering",
      "description":"Mashling gateway recipes from TIBCOSoftware Engineering",
      "publish":[  
        {  
          "recipe":"event-dispatcher-router-mashling",
          "featured":true
        },
        {  
          "recipe":"rest-conditional-gateway",
          "featured":false
        }
      ]
    },
    {  
      "provider":"TIBCOSoftware Services",
      "description":"Mashling gateway recipes from TIBCO Services",
      "publish":[]
    }
  ]
}
```

#### Submitting a New/Updated Recipe
Create a pull request for the recipe to be reviewed and merged into this repository. To publish/unpublish a recipe on https://www.mashling.io/moreRecipes, create a pull request for the updated recipe_registry.json.

#### License
mashling-recipes is licensed under a BSD-type license. See license text [here](https://github.com/TIBCOSoftware/mashling-recipes/blob/master/TIBCO%20LICENSE.txt).

#### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/mashling-recipes/issues).
