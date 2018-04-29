#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2D backBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture2D(backBuffer,fragTexCoords);

	return color;
}

void main()
{
	fragColor = fetchColor();
}