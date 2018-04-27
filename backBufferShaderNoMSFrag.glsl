#version 120

varying vec2 fragTexCoords;

uniform sampler2D backBuffer;

vec4 fetchColor()
{
	return texture2D(backBuffer,fragTexCoords);
}

void main()
{
	gl_FragColor = fetchColor();
	if(gl_FragColor.a < 0.1)
	{
		discard;
	}
}