package gohome

// BackBuffer Shader
const (
	BACKBUFFER_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 150

out vec2 fragTexCoords;

uniform float depth;

vec2 vertices[6];
vec2 texCoords[6];

void setValues()
{
	vertices[0] = vec2(-1.0,-1.0);
	vertices[1] = vec2(1.0,-1.0);
	vertices[2] = vec2(1.0,1.0);
	vertices[3] = vec2(1.0,1.0);
	vertices[4] = vec2(-1.0,1.0);
	vertices[5] = vec2(-1.0,-1.0);

	texCoords[0] = vec2(0.0,0.0);
	texCoords[1] = vec2(1.0,0.0);
	texCoords[2] = vec2(1.0,1.0);
	texCoords[3] = vec2(1.0,1.0);
	texCoords[4] = vec2(0.0,1.0);
	texCoords[5] = vec2(0.0,0.0);
}

void main()
{
	setValues();
	fragTexCoords = texCoords[gl_VertexID];
	gl_Position = vec4(vertices[gl_VertexID],depth,1.0);
}`
	BACKBUFFER_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2DMS BackBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);
	ivec2 texCoords = ivec2(fragTexCoords * textureSize(BackBuffer));

	for(int i = 0;i<8;i++)
	{
		color += texelFetch(BackBuffer,texCoords,i);
	}
	color /= 8.0;

	return color;
}

void main()
{
	fragColor = fetchColor();
	if(fragColor.a < 0.1)
		discard;
}`
	BACKBUFFER_NOMS_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec2 vertex;
attribute vec2 texCoord;

varying vec2 fragTexCoords;

uniform float depth;

void main()
{
	fragTexCoords = texCoord;
	gl_Position = vec4(vertex,depth,1.0);
}`
	BACKBUFFER_NOMS_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying vec2 fragTexCoords;

uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	return texture2D(BackBuffer,fragTexCoords);
}

void main()
{
	gl_FragColor = fetchColor();
	if(gl_FragColor.a < 0.1)
		discard;
}`
)

// Lines3D Shader
const (
	LINES_3D_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec4 color;

varying vec4 fragColor;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragColor = color;
}`
	LINES_3D_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying vec4 fragColor;

void main()
{
    gl_FragColor = fragColor;
}`
)

// PointLight Shadowmap Shader
const (
	POINTLIGHT_SHADOWMAP_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 150

in vec3 vertex;
in vec3 normal;
in vec2 texCoord;
in vec3 tangent;

out vec2 geoTexCoord;

uniform mat4 transformMatrix3D;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	geoTexCoord = texCoord;
}`
	POINTLIGHT_SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 150

in vec2 fragTexCoord;
in vec4 fragPos;

uniform vec3 lightPos;
uniform float farPlane;
uniform struct Material
{
	bool DiffuseTextureLoaded;
	float transparency;
} material;
uniform	sampler2D materialdiffuseTexture;


vec4 fetchColor()
{	
	vec4 col;
	if(material.DiffuseTextureLoaded)
	{
		col = texture2D(materialdiffuseTexture,fragTexCoord);
	}
	else
	{
		col = vec4(1.0,1.0,1.0,1.0);
	}

	col.w *= material.transparency;
}

void main()
{
	vec4 color = fetchColor();
	if(color.a < 0.1)
		discard;
	float lightDistance = length(fragPos.xyz - lightPos);
	lightDistance = lightDistance / farPlane;
	gl_FragDepth = lightDistance;
}`
	POINTLIGHT_SHADOWMAP_SHADER_GEOMETRY_SOURCE_OPENGL string = `
#version 150

layout(triangles) in;
layout(triangle_strip,max_vertices=18) out;

in vec2 geoTexCoord[];

out vec2 fragTexCoord;
out	vec4 fragPos;

uniform mat4 lightSpaceMatrices[6];
uniform mat4 projectionMatrix3D;

void main()
{
	for(int face = 0;face < 6;++face)
	{
		gl_Layer = face;
		for(int i = 0;i<3;++i)
		{
			fragPos = gl_in[i].gl_Position;
			gl_Position = projectionMatrix3D * lightSpaceMatrices[face] * fragPos;
			switch(i)
			{
			    case 0:
			        fragTexCoord = geoTexCoord[0];
			        break;
			    case 1:
                	fragTexCoord = geoTexCoord[1];
                	break;
                case 2:
                    fragTexCoord = geoTexCoord[2];
                	break;
                default:
                    break;
			}
			EmitVertex();
		}
		EndPrimitive();
	}
}`
	POINTLIGHT_SHADOWMAP_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 150

in vec3 vertex;
in vec3 normal;
in vec2 texCoord;
in vec3 tangent;
in mat4 transformMatrix3D;

out	vec2 geoTexCoord;


void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	geoTexCoord = texCoord;
}`
)

// Shadowmap Shader
const (
	SHADOWMAP_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;


varying vec2 fragTexCoord;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}`
	SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying	vec2 fragTexCoord;

uniform struct Material
{
	bool DiffuseTextureLoaded;
	float transparency;
} material;
uniform sampler2D materialdiffuseTexture;

vec4 getDiffuseTexture()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(materialdiffuseTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
}

void main()
{
	vec4 texDifCol = getDiffuseTexture();
	texDifCol.w *= material.transparency;

	if(texDifCol.a < 0.1)
	{
		discard;
	}
}`
	SHADOWMAP_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;
attribute mat4 transformMatrix3D;

varying vec2 fragTexCoord;


uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}`
)
