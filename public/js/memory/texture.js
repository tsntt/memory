const textureLoader = new THREE.TextureLoader();
const baseImage = document.getElementById('cardUvMap')

function dataImage(img) {
  const c = new OffscreenCanvas(img.width, img.height)
  const ctx = c.getContext('2d')
  
  ctx.drawImage(img, 0, 0)

  return ctx.getImageData(0, 0, c.width, c.height)
}

const Material = (() => {
  const dataBaseImage = dataImage(baseImage)

  const materials = {
    BaseMaterial: newMaterial(baseImage.src)
  }

  function newMaterial(src) {
    const texture = textureLoader.load(src)
    texture.flipY = false

    const material = new THREE.MeshPhongMaterial({ map: texture });
    
    return material
  }

  const imgworker = new Worker('../assets/js/memory/worker.js');

  const NewUVMap = data => { imgworker.postMessage(data); };

  imgworker.addEventListener('message', event => {
    materials[event.data.id] = newMaterial( event.data.base64 );
  });

  class Material {
    MakeCustomMaterials(arr) {
      for (let i = 0; i < arr.length; i++) {
        NewUVMap({str: arr[i], baseImage: dataBaseImage.data})
      }
    }

    Get(key) {
      return materials[key]
    }
  }

  return Material
})();