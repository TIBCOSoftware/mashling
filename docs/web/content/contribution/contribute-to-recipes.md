---
title: Contributing to Project Mashling Recipes
weight: 9100
pre: "<i class=\"fa fa-asterisk\" aria-hidden=\"true\"></i> "
---

#### Creating a recipe
A [recipe](https://github.com/TIBCOSoftware/mashling-recipes) can be created by customizing an existing recipe in https://www.mashling.io/moreRecipes or by using an editor such as [VS code](https://code.visualstudio.com/) enabled with a [mashling extension] (https://github.com/TIBCOSoftware/vscode-extension-mashling)

#### Adding a recipe to [mashling-recipes](https://github.com/TIBCOSoftware/mashling-recipes) repo.
A recipe should be contained in its own folder under 'recipes' folder. The recipe folder is composed of a gateway json file, README.md, an optional icon image file, optional Gopkg.lock and Gopkg.toml files. In the absence of the icon image file, the default Mashling icon image is used by mashling.io for the recipe. When the icon image file is present, the Mashling json file should have an icon image file field as follows:

```
{
	"mashling_schema": "1.0",
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

Gopkg.lock and Gopkg.toml contains the specific dependent library versions to be used during the gateway binary compilation. 

#### Publishing a Recipe

[recipe_registry.json](https://github.com/TIBCOSoftware/mashling-recipes/blob/master/recipe_registry.json) contains the list of recipe providers and the recipes to publish. The recipe folder name should be added to the publish field for the recipe to be made available in mashling.io. For example, "event-dispatcher-router-mashling" and "rest-conditional-gateway" recipes binaries are built and made downloadable from mashling.io given the following recipe_registry.json. Setting "featured" to "true" adds the recipe to the featured recipe list in mashling.io.

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
